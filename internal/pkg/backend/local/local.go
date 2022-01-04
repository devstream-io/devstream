package local

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const DefaultStateFile = ".devstream/devstream.state"

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

// Write is used to writes the data to local file.
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
