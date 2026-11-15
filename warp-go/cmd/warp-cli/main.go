package main

import (
	"fmt"
	"os"

	"github.com/robertpelloni/warp/warp-go/pkg/core"
)

func main() {
	fmt.Println("Warp CLI (Go Backend)")
	fmt.Println("---------------------")

	// Initialize the core framework
	engine := core.NewEngine()

	// Parse arguments (stub)
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		fmt.Printf("Executing command: %s\n", cmd)
		engine.Execute(cmd)
	} else {
		fmt.Println("Running in interactive mode...")
		engine.RunInteractive()
	}
}
