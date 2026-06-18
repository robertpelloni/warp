package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Warp Go Initialized")
		fmt.Println("Usage: warp <command> [arguments]")
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "run":
		runCmd := flag.NewFlagSet("run", flag.ExitOnError)
		agent := runCmd.String("agent", "default", "Agent to use")
		runCmd.Parse(os.Args[2:])
		fmt.Printf("Warp Go running agent: %s\n", *agent)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
