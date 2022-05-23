package os

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
)

// ExecInSystem can exec a command with some params in system.
// All logs produced by command would be print to stdout and write into logsBuffer if it is not nil.
// Warning: Don't let other people control the shellCommand Param
// Warning: The ExecInSystem func is insecure when provide as pkg/utils in Framework.
// That will Cause Command injection.
func ExecInSystem(execPath string, shellCommand string, logsBuffer *bytes.Buffer, print bool) error {
	fmt.Printf("Shell Command: %s\n", shellCommand)

	cmd := exec.Command("sh", "-c", shellCommand)
	cmd.Dir = execPath

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// print logs
	var lock sync.Mutex
	outReader := bufio.NewReader(stdout)
	errReader := bufio.NewReader(stderr)
	printLog := func(reader *bufio.Reader, stdType string) {
		for {
			line, err := reader.ReadString('\n')
			if print {
				fmt.Printf("%s: %s", stdType, line)
			}
			if logsBuffer != nil {
				lock.Lock()
				logsBuffer.WriteString(line)
				lock.Unlock()
			}
			if err != nil || err == io.EOF {
				break
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		printLog(outReader, "Stdout")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		printLog(errReader, "Stderr")
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	wg.Wait()
	return cmd.Wait()
}

// SafeExecInSystem can exec a command with some params in system.
// All logs produced by command would be print to stdout and write into logsBuffer if it is not nil
// Warning: The SafeExecInSystem func is not enough secure to defeat the command injection.
// This Func provide a way to avoid params to pollute the whole command.
// Warning: you still can use like command `python -c "{python shell}"` to introduce the code injection.
// even the `kubectl exec` to execute in your k8s.
// when using this func in Plugin, you should take care of parameters for the cmd has no such dangerous functions.
func SafeExecInSystem(execPath string, cmdName string, params []string, logsBuffer *bytes.Buffer, print bool) error {
	fmt.Printf("Exec: %s\n", cmdName)
	fmt.Printf("Params : %s\n", strings.Join(params, " "))
	if IsShell(cmdName) {
		return errors.New("Shell command detected.")
	}
	cmd := exec.Command(cmdName, params...)
	cmd.Dir = execPath

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// print logs
	var lock sync.Mutex
	outReader := bufio.NewReader(stdout)
	errReader := bufio.NewReader(stderr)
	printLog := func(reader *bufio.Reader, stdType string) {
		for {
			line, err := reader.ReadString('\n')
			if print {
				fmt.Printf("%s: %s", stdType, line)
			}
			if logsBuffer != nil {
				lock.Lock()
				logsBuffer.WriteString(line)
				lock.Unlock()
			}
			if err != nil || err == io.EOF {
				break
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		printLog(outReader, "Stdout")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		printLog(errReader, "Stderr")
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	wg.Wait()
	return cmd.Wait()
}

func IsShell(cmdName string) bool {
	if strings.Contains(cmdName, "sh") || strings.Contains(cmdName, "cmd") {
		return true
	} else {
		return false
	}
}
