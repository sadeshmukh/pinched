package tools

type Tool struct {
	Name        string
	Description string
	Parameters  map[string]any
	Required    []string
	Exec        func(params map[string]interface{}) (string, error)
}

var All = []Tool{
	SearchTool,
	CoolifyDeploy,
}
