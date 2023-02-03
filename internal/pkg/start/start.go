package start

import (
	"fmt"
	"time"

	"github.com/devstream-io/devstream/internal/pkg/start/tool"
)

func Start() error {
	fmt.Println("I'll give some tools for you.")
	time.Sleep(time.Second)
	fmt.Println("Are you ready?")
	time.Sleep(time.Second)
	fmt.Println("Let's get started.")
	fmt.Println()

	err := installToolsIfNotExist()
	if err != nil {
		return err
	}

	fmt.Println("\nEverything is going well now.\nEnjoy it!☺️")
	fmt.Println()
	return nil
}

func installToolsIfNotExist() error {
	for _, t := range tool.GetTools() {
		if !t.IfExists() {
			if err := t.Install(); err != nil {
				return err
			}
		}
		if t.IfStopped != nil && t.IfStopped() {
			if err := t.Start(); err != nil {
				return err
			}
		}
		fmt.Printf("✅ %s is ready.\n", t.Name)
	}
	return nil
}
