package logdr

import (
	"strings"
	"testing"
)

func TestFramesStartsAtCaller(t *testing.T) {
	got := frames(0, 3)
	if len(got) == 0 || !strings.Contains(got[0], "logger_test.go") || !strings.Contains(got[0], "TestFramesStartsAtCaller") {
		t.Fatalf("unexpected first frame: %v", got)
	}
}

func helper() []string { return frames(1, 3) }

func TestFramesSkipsOneLevel(t *testing.T) {
	got := helper()
	if len(got) == 0 || !strings.Contains(got[0], "TestFramesSkipsOneLevel") {
		t.Fatalf("expected caller of helper first: %v", got)
	}
}