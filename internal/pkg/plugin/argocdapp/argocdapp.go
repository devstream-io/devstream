package argocdapp

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const defaultYamlPath = "./app.yaml"

const (
	ActionApply  Action = "apply"
	ActionDelete Action = "delete"
)

type Action string

func kubectlAction(action Action, filename string) error {
	cmd := exec.Command("kubectl", string(action), "-f", filename)
	cOut, err := cmd.CombinedOutput()
	if err != nil {
		// TODO(Daniel Hu): Handle the Error below:
		// Error from server (NotFound): error when deleting "./app.yaml": applications.argoproj.io "hello" not found
		log.Printf("failed to exec: < %s >", cmd.String())
		log.Printf("exec logs: < %s >. got error: %s", string(cOut), err)
		return err
	}
	log.Println(strings.TrimSuffix(string(cOut), "\n"))
	return nil
}

func writeContentToTmpFile(file string, content string, param *Param) error {
	t, err := template.New("app").Option("missingkey=error").Parse(content)
	if err != nil {
		return err
	}

	output, err := os.Create(file)
	if err != nil {
		return err
	}

	err = t.Execute(output, param)
	if err != nil {
		if strings.Contains(err.Error(), "can't evaluate field name") {
			msg := err.Error()
			start := strings.Index(msg, "<")
			end := strings.Index(msg, ">")
			return fmt.Errorf("plugin argocdapp needs options%s but it's missing from the config file", msg[start+1:end])
		} else {
			return fmt.Errorf("executing tpl: %s", err)
		}
	}
	return nil
}
