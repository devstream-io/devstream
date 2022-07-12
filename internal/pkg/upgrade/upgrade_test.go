package upgrade

import (
	"testing"
)

func Test_checkForUpdates(t *testing.T) {
	type args struct {
		old string
		new string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "true case with prefix", args: args{old: "v0.6.0", new: "v0.7.0"}, want: true},
		{name: "true case without prefix", args: args{old: "0.6.3", new: "0.7.1"}, want: true},
		{name: "true case with and without prefix", args: args{old: "1.16.100", new: "v1.17.0"}, want: true},

		{name: "false case with prefix", args: args{old: "v0.7.0", new: "v0.7.0"}, want: false},
		{name: "false case without prefix", args: args{old: "0.7.3", new: "0.7.3"}, want: false},
		{name: "true case with and without prefix", args: args{old: "1.18.100", new: "v1.17.0"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkUpgrade(tt.args.old, tt.args.new)
			if err != nil || got != tt.want {
				t.Errorf("checkForUpdates() error = %v, old version = %v, new version = %v, want = %v, got = %v", err, tt.args.old, tt.args.new, tt.want, got)
				return
			}
		})
	}
}
