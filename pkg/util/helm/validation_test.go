package helm

import (
	"reflect"
	"testing"
)

func Test_defaults(t *testing.T) {
	// base
	// timeout := 1h
	tests := []struct {
		name string
		got  HelmParam
		want HelmParam
	}{
		// TODO: Add test cases.
		{"base",
			HelmParam{Repo{"test", ""},
				Chart{
					Timeout: "",
				}},
			HelmParam{
				Repo{"test", ""},
				Chart{
					Timeout: "5m0s",
				}}},
		{"case timeout := 1h",
			HelmParam{Repo{"test", ""},
				Chart{
					Timeout: "1h",
				}},
			HelmParam{
				Repo{"test", ""},
				Chart{
					Timeout: "1h",
				}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaults(&tt.got)
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("validate() = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func Test_validate(t *testing.T) {
	type args struct {
		param *HelmParam
	}
	tests := []struct {
		name string
		args args
		want int // error count
	}{
		// TODO: Add test cases.
		{"base", args{&HelmParam{
			Repo{Name: "argo", URL: "https://argoproj.github.io/argo-helm"},
			Chart{ChartName: "argo/argo-cd"},
		}}, 0},
		{"one required field validation error", args{&HelmParam{
			Repo{Name: "argo", URL: "https://argoproj.github.io/argo-helm"},
			Chart{ChartName: ""},
		}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validate(tt.args.param); len(got) != tt.want {
				t.Logf("got errors' length: %d\n", len(got))
				t.Errorf("validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultsAndValidate(t *testing.T) {
	type args struct {
		param *HelmParam
	}
	type want struct {
		HelmParam
		errCount int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		// TODO: Add test cases.
		{"base", args{&HelmParam{
			Repo{Name: "argo", URL: "https://argoproj.github.io/argo-helm"},
			Chart{ChartName: "argo/argo-cd"},
		}}, want{HelmParam{
			Repo{Name: "argo", URL: "https://argoproj.github.io/argo-helm"},
			Chart{ChartName: "argo/argo-cd", Timeout: "5m0s"},
		}, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultsAndValidate(tt.args.param)
			if len(got) != tt.want.errCount {
				t.Errorf("DefaultsAndValidate(): errorCount = %v, want %v", len(got), tt.want.errCount)
			}
			if !reflect.DeepEqual(*tt.args.param, tt.want.HelmParam) {
				t.Errorf("DefaultsAndValidate(): HelmParam= %v, want %v", *tt.args.param, tt.want.HelmParam)
			}
		})
	}
}
