package main

type Task struct {
	Source   string
	Content  string
	Response chan string
}

func TaskIngestor() chan Task {
	tasks := make(chan Task)

	go func() {
		for task := range tasks {
			tools := []Tool{SearchTool}
			res := aiResponseWithTools(task.Content, tools)
			task.Response <- res
		}
	}()

	return tasks

}
