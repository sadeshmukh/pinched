package main

import (
	"fmt"
	"strings"

	"github.com/sadeshmukh/pinched/secrets"
	"github.com/sadeshmukh/pinched/tools"
)

type Task struct {
	Source   string
	Content  string
	Response chan string
}

func TaskIngestor() chan Task {
	tasks := make(chan Task)

	go func() {
		for task := range tasks {
			if strings.HasPrefix(task.Content, "pinched secret") {
				parts := strings.SplitN(task.Content, " ", 5)
				if len(parts) < 3 {
					task.Response <- "format: `pinched secret [set/get] [name] [value (if set)]`"
				}
				action := parts[2]
				name := parts[3]
				switch action {
				case "set":
					if len(parts) < 5 {
						task.Response <- "format: `pinched secret set [name] [value]`"
						continue
					}
					value := parts[4]
					err := secrets.StoreSecret(name, value)
					if err != nil {
						task.Response <- "error storing secret: " + err.Error()
					} else {
						task.Response <- "secret `" + name + "` stored successfully"
					}

				case "get":
					value, err := secrets.GetSecret(name)
					if err != nil {
						task.Response <- "error getting secret: " + err.Error()
					} else {
						task.Response <- "secret `" + name + "`: " + value
					}
				case "list":
					task.Response <- "TODO: list secrets"

				default:
					task.Response <- "format: `pinched secret [set/get] [name] [value (if set)]`"
				}
				continue
			}

			// metaprompt prepend before direct task.content
			prompt := fmt.Sprintf(`You are Pinched, a personal assistant for Sahil, provided certain tools to manage their infrastructure. Keep in mind the formatting of the source: %s. When given secrets in the form {{SECRET}}, repeat them verbatim and they will be substituted later on. Do not use tools unless necessary. User prompt: %s`, task.Source, task.Content)

			res := aiResponseWithTools(prompt, tools.All)
			task.Response <- res
		}
	}()

	return tasks

}
