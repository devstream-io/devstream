package create

import (
	"fmt"
	"time"
)

type Param struct {
	GithubUsername    string
	GithubToken       string
	DockerhubUsername string
	DockerhubToken    string
	Language          string
	Framework         string
}

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

	params, err := getParams()
	if err != nil {
		return err
	}

	return create(params)
}

// TODO: @jf
func create(params *Param) error {
	err := createRepo(params)
	if err != nil {
		return err
	}

	return createApp(params)
}

// TODO(daniel-hutao): support python/flask first
func createRepo(params *Param) error {
	fmt.Printf("Lang: %s, Fram: %s\n", params.Language, params.Framework)
	return nil
}

func createApp(params *Param) error {
	return nil
}
