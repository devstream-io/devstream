package types

import (
	"errors"
	"fmt"
)

type backendOptionErr struct {
	backendType Type
}

func NewBackendOptionErr(typ Type) error {
	return backendOptionErr{backendType: typ}
}

func (err backendOptionErr) Error() string {
	return fmt.Sprintf("invaid %s sate config", err.backendType)
}

func IsBackendOptionErr(err error) bool {
	return errors.As(err, &backendOptionErr{})
}

type invalidBackendErr struct {
	backendType string
}

func NewInvalidBackendErr(typ string) error {
	return invalidBackendErr{backendType: typ}
}

func (err invalidBackendErr) Error() string {
	return fmt.Sprintf("the backend type < %s > is illegal", err.backendType)
}

func IsInvalidBackendErr(err error) bool {
	return errors.As(err, &invalidBackendErr{})
}
