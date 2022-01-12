package os

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
)

// ExecInSystem can exec a command with some params in system.
// All logs produced by command would be print to stdout and write into logsBuffer if it is not nil
func ExecInSystem(execPath string, params []string, logsBuffer *bytes.Buffer, print bool) error {
	paramStr := strings.Join(params, " ")
	fmt.Printf("Params : %s\n", paramStr)

	c := "-c"
	cmdName := "sh"

	cmd := exec.Command(cmdName, c, paramStr)
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
