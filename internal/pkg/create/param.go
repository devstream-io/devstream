package create

import (
	"fmt"
	"time"
)

// TODO: @bird
func getParams() (*Param, error) {
	lang, err := getLanguage()
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second)
	fmt.Println("\nPlease choose a framework next.")
	time.Sleep(time.Second)

	fram, err := getFramework()
	if err != nil {
		return nil, err
	}
	return &Param{
		Language:  lang,
		Framework: fram,
	}, nil
}
