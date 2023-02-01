package create

import (
	"fmt"
	"time"

	"github.com/devstream-io/devstream/internal/pkg/create/param"
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

	params, err := param.GetParams()
	if err != nil {
		return err
	}

	return create(params)
}

// TODO: @jf
func create(params *param.Param) error {
	err := createRepo(params)
	if err != nil {
		return err
	}

	return createApp(params)
}

// TODO(daniel-hutao): support python/flask first
func createRepo(params *param.Param) error {
	fmt.Printf("Lang: %s, Fram: %s\n", params.Language, params.Framework)
	return nil
}

func createApp(params *param.Param) error {
	return nil
}
