package main


func TaskIngestor() chan Task {
	tasks := make(chan Task)


	go func() {
		for task := range tasks {
			// holup I'll figure this out
		}
	}()

	return tasks

}

