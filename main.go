package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /query/{q}", func (w http.ResponseWriter, r *http.Request) {
		q := r.PathValue("q")
		aiResponse(q)
		fmt.Print("got test")
		w.Write([]byte("successful"))
	})

	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		fmt.Println(err)
	}
}



