package main

import (
	"fmt"
	"strings"
	"time"

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
					continue
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

			date := fmt.Sprintf("Current date and time: %s", time.Now().Local().Format("2006-01-02"))
			// metaprompt prepend before direct task.content
			slack_prompt := "Strictly adhere to this formatting: slack uses mrkdwn for formatting, so you can use *bold*, _italic_, and `code` formatting in your responses, as well as line breaks, block quotes, code blocks, and lists. For links, use this formatting: <http://www.example.com|link>."

			prompt := fmt.Sprintf(`You are Pinched, a personal assistant for Sahil, provided certain tools to manage their infrastructure. Keep in mind the formatting of the source: %s. When given secrets in the form {{SECRET}}, repeat them verbatim and they will be substituted later on. Do not use tools unless necessary. Generally, use the location WDUB if none other is provided. %s`,
				task.Source,
				date,
			)

			if task.Source == "slack" {
				prompt = prompt + "\n\n" + slack_prompt
			}

			prompt = prompt + "\n\n" + task.Content

			exitLoopTool := tools.Tool{
				Name:        "end_with_resp",
				Description: "Call this when you've finished the task and are ready to respond with a summary of your actions taken, if applicable, or a final answer. This can be used to the exit even if no tools have been called.",
				Parameters: map[string]any{
					"response": map[string]any{
						"type":        "string",
						"description": "The final response to the user's query. Can be a summary of actions taken, if applicable, or a final answer.",
					},
				},
				Required: []string{"response"},
				Exec: func(params map[string]interface{}) (string, error) {
					resp := params["response"].(string)

					return "|END|" + resp, nil
				},
			}

			res := aiResponseWithTools(prompt, append(tools.All, exitLoopTool))
			task.Response <- res
		}
	}()

	return tasks

}
