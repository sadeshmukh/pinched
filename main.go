package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// see coolify.go for why this doesn't work, tl;dr coolify API just doesn't work at all

	// _res, err := tools.CoolifyDeploy.Exec(
	// 	map[string]interface{}{
	// 		"repo": "sadeshmukh/listen",
	// 		// "url":  "https://listen.halceon.dev",
	// 	},
	// )
	// if err != nil {
	// 	fmt.Printf("deploy error: %v\n", err)
	// } else {
	// 	fmt.Printf("deploy response: %s\n", _res)
	// }

	// res, err := tools.BARTTripTestLookupTool.Exec(nil)
	// if err != nil {
	// 	fmt.Printf("BART real-time error: %v\n", err)
	// } else {
	// 	fmt.Printf("BART real-time response: %s\n", res)
	// }

	// res, err := tools.BARTStationTool.Exec(nil)
	// if err != nil {
	// 	fmt.Printf("BART station list error: %v\n", err)
	// } else {
	// 	fmt.Printf("BART station list response: %s\n", res)
	// }

	// select {}

	tasks := TaskIngestor()
	go func() {
		if err := Discord(tasks); err != nil {
			panic(err)
		}
	}()
	err = Slack(tasks)
	if err != nil {
		panic(err)
	}

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
