package tools

import (
	"bytes"
	"strings"

	"golang.org/x/crypto/ssh"
)

// will also copy the newt config in a compose override, then tunnel w/ pangolin to preferred URL eventually

var DeployTool = Tool{
	Name:        "deploy",
	Description: "Deploy a public git repository, Docker-compose based, to a manually-configured server",
	Parameters: map[string]any{
		"repo": map[string]any{
			"type":        "string",
			"description": "The Github repo to deploy (e.g. sadeshmukh/pinched)",
		},
		"compose_path": map[string]any{
			"type":        "string",
			"description": "Path to docker-compose file in the repo (defaults to docker-compose.yml)",
		},
		"server": map[string]any{
			"type":        "string",
			"description": "Server name to deploy to (must be accessible over tailscale). Defaults to `sahil@nest`.",
		},
		"url": map[string]any{
			"type":        "string",
			"description": "Optional URL to expose (must exist in Pangolin already)",
		},
	},
	Exec: func(params map[string]interface{}) (string, error) {
		// repo, _ := params["repo"].(string)
		composePath, _ := params["compose_path"].(string)
		server, _ := params["server"].(string)

		if composePath == "" {
			composePath = "docker-compose.yml"
		}
		if server == "" {
			server = "sahil@nest"
		}

		// must be accessible over tailscale magic dns at that server hostname

		sshConfig := &ssh.ClientConfig{
			User:            strings.Split(server, "@")[0],
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		client, err := ssh.Dial("tcp", strings.Split(server, "@")[1]+":22", sshConfig)
		if err != nil {
			return "", err
		}
		defer client.Close()

		session, err := client.NewSession()
		if err != nil {
			return "", err
		}
		defer session.Close()

		var output bytes.Buffer
		session.Stdout = &output
		session.Stderr = &output

		// TODO: actually do something
		return "", nil
	},
}
