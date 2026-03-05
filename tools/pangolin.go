package tools

var PangolinInfoTool = Tool{
	Name:        "pangolin_info",
	Description: "Gets various info about Pangolin",
	Parameters: map[string]any{
		"type": map[string]any{
			"type":        "string",
			"description": "One of: resources, subdomains, sites",
		},
	},
}

var PangolinTunnelTool = Tool{
	Name:        "pangolin_tunnel",
	Description: "Creates tunnels with Pangolin.",
	Parameters: map[string]any{
		"action": map[string]any{
			"type":        "string",
			"description": "Action to perform: create, list",
		},
		"subdomain": map[string]any{
			"type":        "string",
			"description": "Subdomain to use for tunnel (must be added already)",
		},
		"destination": map[string]any{
			"type":        "string",
			"description": "Destination address for the tunnel, e.g. site:http://localhost:4321",
		},
	},
}
