package tools

import (
	"os"
	"testing"
)

func TestReadFilesTool(t *testing.T) {
	tool := &ReadFilesTool{}
	tmpFile, _ := os.CreateTemp("", "testfile")
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("hello world")
	tmpFile.Close()

	args := map[string]interface{}{"path": tmpFile.Name()}
	content, err := tool.Execute(args)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if content != "hello world" {
		t.Errorf("Expected 'hello world', got %q", content)
	}
}

func TestExecuteCommandTool(t *testing.T) {
	tool := &ExecuteCommandTool{}
	args := map[string]interface{}{"command": "echo 'test'"}
	output, err := tool.Execute(args)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if output != "test\n" {
		t.Errorf("Expected 'test\\n', got %q", output)
	}
}
