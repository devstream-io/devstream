package response

import (
	"encoding/json"
	"fmt"
	"github.com/devstream-io/devstream/internal/log"

	"gopkg.in/yaml.v3"
)

type StatusCode int
type MessageText string

type Response struct {
	Status  StatusCode  `json:"status" yaml:"status"`
	Message MessageText `json:"message" yaml:"message"`
	Log     string      `json:"log" yaml:"log"`
}

var (
	StatusOK    StatusCode = 0
	StatusError StatusCode = 1
)

var (
	MessageOK    MessageText = "OK"
	MessageError MessageText = "ERROR"
)

func New(status StatusCode, message MessageText, log string) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Log:     log,
	}
}

func (r *Response) Print(format string) {
	log.Debugf("Format: %s", format)
	switch format {
	case "json":
		r.printJSON()
	case "yaml":
		r.printYAML()
	default:
		r.printRaw()
	}
}

func (r *Response) printRaw() {
	fmt.Println(r.toRaw())
}

func (r *Response) printJSON() {
	fmt.Println(r.toJSON())
}

func (r *Response) printYAML() {
	fmt.Println(r.toYAML())
}

func (r *Response) toRaw() string {
	return r.Log
}

func (r *Response) toJSON() string {
	str, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(str)
}

func (r *Response) toYAML() string {
	str, err := yaml.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(str)
}
