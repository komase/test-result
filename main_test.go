package main

import (
	"flag"
	"os"
	"testing"
	"time"
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

func Test_GroupResultsByTestStatus(t *testing.T) {
	now := time.Now()
	testName1 := "main"
	testName2 := "bar"
	testName3 := "foo"
	output1 := "=== RUN   Test_main\\n"
	output2 := "--- PASS: Test_main (0.00s)\\n"
	output3 := "=== RUN   Test_bar\\n"
	output4 := "bar\\n\""
	output5 := "    bar_test.go:9: bar()\\n"
	output6 := "--- FAIL: Test_bar (0.00s)\\n"
	output7 := "=== RUN   Test_foo2\\n"
	output8 := "    foo_test.go:13: \\n"
	output9 := "--- SKIP: Test_foo2 (0.00s)\\n"
	elapsed := 0.1

	results := []Result{
		{Time: now, Action: "start", Package: "pkg1"},
		{Time: now, Action: "run", Package: "pkg1", Test: &testName1},
		{Time: now, Action: "output", Package: "pkg1", Test: &testName1, Output: &output1},
		{Time: now, Action: "output", Package: "pkg1", Test: &testName1, Output: &output2},
		{Time: now, Action: "pass", Package: "pkg1", Test: &testName1, Elapsed: &elapsed},
		{Time: now, Action: "run", Package: "pkg2", Test: &testName2},
		{Time: now, Action: "output", Package: "pkg2", Test: &testName2, Output: &output3},
		{Time: now, Action: "output", Package: "pkg2", Test: &testName2, Output: &output4},
		{Time: now, Action: "output", Package: "pkg2", Test: &testName2, Output: &output5},
		{Time: now, Action: "output", Package: "pkg2", Test: &testName2, Output: &output6},
		{Time: now, Action: "fail", Package: "pkg2", Test: &testName2, Elapsed: &elapsed},
		{Time: now, Action: "run", Package: "pkg3", Test: &testName3, Elapsed: &elapsed},
		{Time: now, Action: "output", Package: "pkg3", Test: &testName3, Output: &output7},
		{Time: now, Action: "output", Package: "pkg3", Test: &testName3, Output: &output8},
		{Time: now, Action: "output", Package: "pkg3", Test: &testName3, Output: &output9},
		{Time: now, Action: "skip", Package: "pkg3", Test: &testName3, Elapsed: &elapsed},
	}

	passResults, failResults, skipResults := groupResultsByTestStatus(results)

	if _, ok := passResults[TestName(testName1)]; !ok {
		t.Errorf("Expected test %s to pass", testName1)
	}

	if _, ok := failResults[TestName(testName2)]; !ok {
		t.Errorf("Expected test %s to fail", testName2)
	}

	if _, ok := skipResults[TestName(testName3)]; !ok {
		t.Errorf("Expected test %s to fail", testName3)
	}

	wantPass := 2
	wantFail := 4
	wantSkip := 3
	gotPass := len(passResults[TestName(testName1)])
	gotFail := len(failResults[TestName(testName2)])
	gotSkip := len(skipResults[TestName(testName3)])

	if gotPass != wantPass || gotFail != wantFail || gotSkip != wantSkip {
		t.Errorf("Unexpected number of results. Got %d pass, %d fail, %d skip. Want %d pass, %d fail, %d skip.",
			gotPass, gotFail, gotSkip, wantPass, wantFail, wantSkip)
	}

}
