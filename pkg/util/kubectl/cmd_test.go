package kubectl

import (
	"testing"
)

var fileName = "nginx-deployment.yaml"

func TestKubeCreate(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{"base", fileName, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := KubeCreate(tt.filename); (err != nil) != tt.wantErr {
				t.Errorf("KubeCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKubeApply(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{"base", fileName, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := KubeApply(tt.filename); (err != nil) != tt.wantErr {
				t.Errorf("KubeApply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKubeDelete(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{"base", fileName, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := KubeDelete(tt.filename); (err != nil) != tt.wantErr {
				t.Errorf("KubeDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_kubectlAction(t *testing.T) {
	tests := []struct {
		name     string
		action   string
		filename string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{"base", APPLY, fileName, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := kubectlAction(tt.action, tt.filename); (err != nil) != tt.wantErr {
				t.Errorf("kubectlAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
