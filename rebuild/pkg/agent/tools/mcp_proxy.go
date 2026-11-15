package tools

import (
	"context"
	"github.com/robertpelloni/warp-rebuild/pkg/mcp"
)

type MCPToolProxy struct {
	client *mcp.Client
	tool   mcp.Tool
}

func (t *MCPToolProxy) Name() string        { return t.tool.Name }
func (t *MCPToolProxy) Description() string { return t.tool.Description }
func (t *MCPToolProxy) Parameters() map[string]interface{} {
	return t.tool.InputSchema
}
func (t *MCPToolProxy) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	return t.client.CallTool(ctx, t.tool.Name, args)
}

func NewMCPToolProxies(client *mcp.Client) ([]*MCPToolProxy, error) {
	tools, err := client.ListTools(context.Background())
	if err != nil {
		return nil, err
	}
	var proxies []*MCPToolProxy
	for _, tool := range tools {
		proxies = append(proxies, &MCPToolProxy{client: client, tool: tool})
	}
	return proxies, nil
}
