package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Task struct {
	Source string
	Content string
	Response chan string
}


func main() {
	err := godotenv.Load()
	
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	tasks := TaskIngestor()
	Discord(tasks)


	// mux := http.NewServeMux()

	// mux.HandleFunc("GET /query/{q}", func (w http.ResponseWriter, r *http.Request) {
	// 	q := r.PathValue("q")
	// 	aiResponse(q)
	// 	fmt.Print("got test")
	// 	w.Write([]byte("successful"))
	// })

	// err = http.ListenAndServe(":3000", mux)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}



