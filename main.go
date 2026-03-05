package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sadeshmukh/pinched/secrets"
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

	err = secrets.StoreSecret("test", "highly highly secret value over here")
	if err != nil {
		panic(err)
	}
	val, err := secrets.GetSecret("test")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	subbed := secrets.Substitute("secret value: {{test}}")
	fmt.Println(subbed)

	tasks := TaskIngestor()
	err = Discord(tasks)
	if err != nil {
		panic(err)
	}

	select {}

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
