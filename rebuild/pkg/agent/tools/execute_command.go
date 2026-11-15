package tools

import (
	"context"
	"fmt"
	"os/exec"
)

type ExecuteCommandTool struct{}

func (t *ExecuteCommandTool) Name() string        { return "execute_command" }
func (t *ExecuteCommandTool) Description() string { return "Executes a shell command." }
func (t *ExecuteCommandTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]interface{}{
				"type":        "string",
				"description": "The command to execute.",
			},
		},
		"required": []string{"command"},
	}
}
func (t *ExecuteCommandTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	command, ok := args["command"].(string)
	if !ok {
		return "", fmt.Errorf("missing command parameter")
	}
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}
