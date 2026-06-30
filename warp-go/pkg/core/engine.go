package core

import "fmt"

type Engine struct {
	State string
}

func NewEngine() *Engine {
	return &Engine{State: "Initialized"}
}

func (e *Engine) Execute(command string) {
	fmt.Printf("Engine executing: %s\n", command)
}

func (e *Engine) RunInteractive() {
	fmt.Println("Engine interactive loop started.")
}
