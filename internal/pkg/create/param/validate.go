package param

import (
	"errors"
	"strings"
)

var validate = func(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("input can not be empty")
	}
	return nil
}
