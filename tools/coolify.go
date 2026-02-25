package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// highly privileged token
// eventually, idea is to also add some sort of autocodegen to automatically create repo & deploy that to coolify
// for now, it'll just be able to use the github app token to deploy existing repo

var CoolifyDeploy = Tool{
	Name:        "coolify.deploy",
	Description: "Deploys a Github repo to Coolify. Must be docker compose based.",
	Parameters: map[string]any{
		"repo": map[string]any{
			"type":        "string",
			"description": "The Github repo to deploy (e.g. sadeshmukh/pinched)",
		},
		"compose_path": map[string]any{
			"type":        "string",
			"description": "Defaults to docker-compose.yml.",
		},
		"url": map[string]any{
			"type":        "string",
			"description": "Optional URL to expose (e.g. https://app.example.com).",
		},
	},
	Required: []string{"repo"},
	Exec: func(params map[string]interface{}) (string, error) {
		coolifyEndpoint := os.Getenv("COOLIFY_ENDPOINT")
		projectID := os.Getenv("COOLIFY_PROJECT_ID")
		environmentID := os.Getenv("COOLIFY_ENVIRONMENT_ID")
		environmentName := os.Getenv("COOLIFY_ENVIRONMENT_NAME")
		serverID := os.Getenv("COOLIFY_SERVER_ID")
		destinationID := os.Getenv("COOLIFY_DESTINATION_ID")
		githubAppID := os.Getenv("COOLIFY_GITHUB_APP_ID")
		coolifyToken := os.Getenv("COOLIFY_TOKEN")

		fixed := []string{
			"project_uuid",
			"environment_name",
			"environment_uuid",
			"server_uuid",
			"destination_uuid",
			"github_app_uuid",
		}

		agent_required := []string{
			"git_repository",
			"git_branch",
			"name",
			"docker_compose_location",
		}

		agent_optional := map[string]any{
			"description":             "string",
			"instant_deploy":          false,
			"force_domain_override":   false,
			"urls":                    []map[string]any{{"name": "string", "url": "string"}},
			"docker_compose_location": "string",
		}

		create_schema := map[string]any{
			"endpoint": coolifyEndpoint,
			"fixed_values": map[string]any{
				"project_uuid":     projectID,
				"environment_uuid": environmentID,
				"environment_name": environmentName,
				"server_uuid":      serverID,
				"destination_uuid": destinationID,
				"github_app_uuid":  githubAppID,
			},
			"fixed_fields":   fixed,
			"agent_required": agent_required,
			"agent_optional": agent_optional,
			"fields": map[string]any{
				"project_uuid":            "string",
				"environment_name":        "string",
				"environment_uuid":        "string",
				"server_uuid":             "string",
				"destination_uuid":        "string",
				"github_app_uuid":         "string",
				"git_repository":          "string",
				"git_branch":              "string",
				"name":                    "string",
				"description":             agent_optional["description"],
				"instant_deploy":          agent_optional["instant_deploy"],
				"force_domain_override":   agent_optional["force_domain_override"],
				"urls":                    agent_optional["urls"],
				"docker_compose_location": agent_optional["docker_compose_location"],
			},
		}
		println(create_schema)

		repo, _ := params["repo"].(string)
		composePath, _ := params["compose_path"].(string)
		urlValue, _ := params["url"].(string)

		if composePath == "" {
			composePath = "docker-compose.yml"
		}

		payload := map[string]any{
			"project_uuid":            projectID,
			"environment_uuid":        environmentID,
			"environment_name":        environmentName,
			"server_uuid":             serverID,
			"destination_uuid":        destinationID,
			"github_app_uuid":         githubAppID,
			"git_repository":          repo,
			"git_branch":              "main",
			"name":                    repo,
			"docker_compose_location": composePath,
			"instant_deploy":          true,
		}

		if urlValue != "" {
			payload["urls"] = []map[string]any{{"name": "app", "url": urlValue}}
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return "", fmt.Errorf("failed to marshal: %w", err)
		}

		endpoint := strings.TrimRight(coolifyEndpoint, "/") + "/applications/private-github-app"
		req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))
		if err != nil {
			return "", fmt.Errorf("failed request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		if coolifyToken != "" {
			req.Header.Set("Authorization", "Bearer "+coolifyToken)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed request: %w", err)
		}
		defer resp.Body.Close()

		respBody := new(bytes.Buffer)
		_, _ = respBody.ReadFrom(resp.Body)
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return "", fmt.Errorf("coolify error (%d): %s", resp.StatusCode, respBody.String())
		}

		return respBody.String(), nil

	},
}

var CoolifyTool = Tool{
	Name:        "coolify",
	Description: "Read, write, deploy to Coolify.",
	Parameters: map[string]any{
		"action": map[string]any{
			"type":        "string",
			"description": "Action to perform: read, deploy",
		},
	},

	Required: []string{"action"},
}
