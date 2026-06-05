package tools

import (
	"fmt"
	"os"
	"os/exec"
)

// Tool represents a capability the agent can use.
type Tool interface {
	Execute(args map[string]interface{}) (string, error)
	GetDefinition() map[string]interface{}
}

// ReadFilesTool reads the content of a file.
type ReadFilesTool struct{}

func (t *ReadFilesTool) GetDefinition() map[string]interface{} {
	return map[string]interface{}{
		"name":        "read_file",
		"description": "Read the content of a file at the given path.",
		"parameters": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{"type": "string"},
			},
			"required": []string{"path"},
		},
	}
}

func (t *ReadFilesTool) Execute(args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		return "", fmt.Errorf("missing path parameter")
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ExecuteCommandTool runs a shell command.
type ExecuteCommandTool struct{}

func (t *ExecuteCommandTool) GetDefinition() map[string]interface{} {
	return map[string]interface{}{
		"name":        "execute_command",
		"description": "Run a shell command and return the output.",
		"parameters": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"command": map[string]interface{}{"type": "string"},
			},
			"required": []string{"command"},
		},
	}
}

func (t *ExecuteCommandTool) Execute(args map[string]interface{}) (string, error) {
	command, ok := args["command"].(string)
	if !ok {
		return "", fmt.Errorf("missing command parameter")
	}
	out, err := exec.Command("sh", "-c", command).CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}
