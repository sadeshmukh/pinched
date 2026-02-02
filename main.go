package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /test", func (w http.ResponseWriter, r *http.Request) {
		fmt.Print("got test")
		w.Write([]byte("successful"))
	})

	http.ListenAndServe(":3000", mux)
}



