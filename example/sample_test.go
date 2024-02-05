//go:build example
// +build example

package sample

import (
	"testing"
)

func Test_bar(t *testing.T) {
	t.Errorf("Test_bar")
}

func Test_bar2(t *testing.T) {
	t.Skip()
}

func Test_bar3(t *testing.T) {
	t.Parallel()
}
