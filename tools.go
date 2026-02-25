package main

type Tool struct {
	Name        string
	Description string
	Parameters  map[string]any
	Required    []string
	Exec        func(params map[string]interface{}) (string, error)
}

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
		// query := params["query"].(string)
		// TODO: actually do something
		return "res", nil
	},
}
