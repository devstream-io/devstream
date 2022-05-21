package local

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/devstream-io/devstream/pkg/util/log"
)

const DefaultStateFile = "devstream.state"

// Local is a default implement for backend.Backend
type Local struct {
	mu       sync.Mutex
	filename string
}

// NewLocal will use DefaultStateFile as statemanager file if filename is not given.
func NewLocal(filename string) *Local {
	var lFile = filename
	if filename == "" {
		lFile = DefaultStateFile
	}

	if _, err := os.Stat(lFile); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(lFile)
		if err != nil {
			log.Fatalf("Creating state file %s failed.", lFile)
		}
		log.Debugf("The state file %s have been created.", lFile)
		defer file.Close()
	}

	return &Local{
		filename: lFile,
	}
}

// Read is used to retrieve the data from local file.
func (l *Local) Read() ([]byte, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := ioutil.ReadFile(l.filename)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Write is used to write the data to local file.
func (l *Local) Write(data []byte) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(l.filename), 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(l.filename, data, 0644); err != nil {
		return err
	}
	return nil
}
