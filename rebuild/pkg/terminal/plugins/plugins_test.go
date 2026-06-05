package plugins

import "testing"

type MockPlugin struct {
	name string
}

func (p *MockPlugin) Name() string { return p.name }
func (p *MockPlugin) Init() error { return nil }

func TestPluginManager(t *testing.T) {
	mgr := NewPluginManager()
	p := &MockPlugin{name: "test-plugin"}

	if err := mgr.LoadPlugin(p); err != nil {
		t.Fatalf("LoadPlugin failed: %v", err)
	}

	plugins := mgr.GetPlugins()
	if len(plugins) != 1 {
		t.Errorf("Expected 1 plugin, got %d", len(plugins))
	}

	if plugins[0].Name() != "test-plugin" {
		t.Errorf("Expected plugin name 'test-plugin', got %s", plugins[0].Name())
	}
}
