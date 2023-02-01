package create

import (
	"fmt"
	"time"
)

func Create() error {
	helloMsg := func() {
		fmt.Println("I'll scaffold a new repository for you.")
		time.Sleep(time.Second)
		fmt.Println("Are you ready?")
		time.Sleep(time.Second)
		fmt.Println("Let's get started.")
		time.Sleep(time.Second)
	}
	helloMsg()

	lang, err := getLanguage()
	if err != nil {
		return err
	}

	time.Sleep(time.Second)
	fmt.Println("\nPlease choose a framework next.")
	time.Sleep(time.Second)

	fram, err := getFramework()
	if err != nil {
		return err
	}

	return createRepo(lang, fram)
	// TODO(daniel-hutao): cicd
}

// TODO(daniel-hutao): support python/flask first
func createRepo(lang, fram string) error {
	fmt.Printf("Lang: %s, Fram: %s\n", lang, fram)
	return nil
}
