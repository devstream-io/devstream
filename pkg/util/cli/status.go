package cli

import (
	"bytes"
	"fmt"
	"time"

	"github.com/briandowns/spinner"

	"github.com/devstream-io/devstream/pkg/util/log"
)

type Status struct {
	spinner       *spinner.Spinner
	logBuffer     *bytes.Buffer
	statusMessage string
	successFormat string
	failureFormat string
}

func StatusForPlugin() *Status {
	spinner := spinner.New(spinner.CharSets[7], 5000*time.Microsecond)
	_ = spinner.Color("yellow")
	return &Status{
		spinner:       spinner,
		logBuffer:     new(bytes.Buffer),
		successFormat: "\x1b[32m✓\x1b[0m %s\n",
		failureFormat: "\x1b[31m✗\x1b[0m %s\n",
	}
}

func (s *Status) Start(status string) {
	s.statusMessage = status
	s.spinner.Suffix = fmt.Sprintf(" %s", status)
	log.RedirectOutput(s.logBuffer)
	s.spinner.Start()
}

// End completes the current status, ending any previous spinning and
// marking the status as success or failure
func (s *Status) End(err error) {
	log.RecoverOutput()
	if err == nil {
		s.spinner.FinalMSG = fmt.Sprintf(s.successFormat, s.statusMessage)
	} else {
		s.spinner.FinalMSG = fmt.Sprintf(s.failureFormat, s.statusMessage)
	}
	// logBuffer contains log during status.Start to status.End
	// we should write this logBuffer to something for error trace
	s.logBuffer.Reset()
	s.spinner.Stop()
}
