package argocdapp

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

func writeContentToTmpFile(file string, content string, param *Param) {
	t, err := template.New("app").Option("missingkey=error").Parse(content)
	if err != nil {
		log.Fatalf("Parse template: %s", err.Error())
	}

	output, err := os.Create(file)
	if err != nil {
		log.Fatalf("Create outputFile: %s", err)
	}

	err = t.Execute(output, param)
	if err != nil {
		if strings.Contains(err.Error(), "can't evaluate field name") {
			msg := err.Error()
			start := strings.Index(msg, "<")
			end := strings.Index(msg, ">")
			log.Fatalf("plugin argocdapp needs options%s but it's missing from the config file", msg[start+1:end])
		} else {
			log.Fatalf("Executing tpl: %s", err)
		}
	}
}

func kubectlApply(file string) {
	cmd := exec.Command("kubectl", "apply", "-f", file)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println(string(strings.TrimSuffix(string(stdout), "\n")))
}
