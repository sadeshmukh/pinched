package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// so this straight up doesn't work because coolify has had this bug for probably over a year that makes this API endpoint completely nonfunctional
// how nobody's ever seen this entirely baffles me, and then I remember nobody actually has ever used the Coolify API for anything other than deploy webhooks
// https://github.com/coollabsio/coolify/issues/5467

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

			"build_pack":    "dockercompose",
			"ports_exposes": "3000",
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

var CoolifyDeployV2 = Tool{
	Name:        "coolify.deploy_v2",
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
		return "", nil
		// time to playwright this entire thing hell yeah
	},
}
