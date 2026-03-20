package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

var SearchTool = Tool{
	Name:        "search_web",
	Description: "Searches the web",
	Parameters: map[string]any{
		"query": map[string]any{
			"type":        "string",
			"description": "The search query",
		},
	},
	Required: []string{"query"},

	Exec: func(params map[string]interface{}) (string, error) {
		fmt.Println("search: querying " + params["query"].(string))
		searchURL := "https://search.hackclub.com/res/v1/web/search" + "?q=" + url.QueryEscape(params["query"].(string))

		request, err := http.NewRequest("GET", searchURL, nil)
		if err != nil {
			return "", err
		}

		request.Header.Set("accept", "application/json")
		request.Header.Set("X-Subscription-Token", os.Getenv("SEARCH_TOKEN"))
		// request.Header.Set("accept-encoding", "gzip")

		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		var result map[string]any
		err = json.Unmarshal(body, &result)
		if err != nil {
			return "", err
		}

		web, ok := result["web"].(map[string]any)
		if !ok || web == nil {
			return "No results found.", nil
		}
		results, ok := web["results"].([]any)
		if !ok || results == nil {
			return "No results found.", nil
		}
		output := ""
		for i, r := range results {
			if i >= 5 {
				break
			}
			sres := r.(map[string]any)
			output += sres["title"].(string) + "\n"
			output += sres["url"].(string) + "\n"
			output += sres["description"].(string) + "\n\n"
		}
		return output, nil

	},
}
