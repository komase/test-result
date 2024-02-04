package main

import (
	"flag"
	"os"
	"testing"
)

func Test_Run(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want int
	}{
		{
			name: "Test_Run",
			args: []string{"-f", "test/results_sample"},
			want: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			commandLineFlagSet = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			if got := run(tt.args); got != tt.want {
				t.Errorf("run() = %v, want %v", got, tt.want)
			}
		})
	}
}
