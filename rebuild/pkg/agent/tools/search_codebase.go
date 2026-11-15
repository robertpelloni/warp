package tools

import (
	"context"
	"fmt"
	"os/exec"
)

type SearchCodebaseTool struct{}

func (t *SearchCodebaseTool) Name() string        { return "search_codebase" }
func (t *SearchCodebaseTool) Description() string { return "Searches the codebase for a pattern." }
func (t *SearchCodebaseTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "The search pattern.",
			},
		},
		"required": []string{"query"},
	}
}
func (t *SearchCodebaseTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("missing query parameter")
	}
	// Use ripgrep-style search for initial parity
	cmd := exec.CommandContext(ctx, "grep", "-r", query, ".")
	out, _ := cmd.CombinedOutput()
	return string(out), nil
}
