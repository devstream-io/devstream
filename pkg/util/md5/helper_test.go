package md5

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCalcFileMD5(t *testing.T) {
	testDir := t.TempDir()
	bFile, _ := os.Create(filepath.Join(testDir, "b.txt"))
	cFile, _ := os.Create(filepath.Join(testDir, "c.txt"))
	_, err := cFile.WriteString("test")
	if err != nil {
		t.Error(err)
	}
	defer bFile.Close()
	defer cFile.Close()
	tests := []struct {
		name     string
		filename string
		want     string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{"base not exist", "a.txt", "", true},
		{"base empty file", bFile.Name(), "d41d8cd98f00b204e9800998ecf8427e", false},
		{"base contented file", cFile.Name(), "098f6bcd4621d373cade4e832627b4f6", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalcFileMD5(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalcFileMD5() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalcFileMD5() = %v, want %v", got, tt.want)
			}
		})
	}
}
