package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Regression test for Bug #5: verify no legacy root main.go exists that would conflict with cmd/ entry points.
func TestNoRootMainGo(t *testing.T) {
	// Walk up from the test file to find the project root (where go.mod is)
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root")
		}
		dir = parent
	}

	// Verify no root main.go exists
	rootMain := filepath.Join(dir, "main.go")
	if _, err := os.Stat(rootMain); err == nil {
		t.Fatal("root main.go still exists — Bug #5 regression: must be removed to avoid build conflict with cmd/ entry points")
	}

	// Verify that the API server builds without conflict
	build := exec.Command("go", "build", "./cmd/api")
	build.Dir = dir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("cmd/api build failed: %v\n%s", err, out)
	}

	// Verify that the worker builds without conflict
	build = exec.Command("go", "build", "./cmd/worker")
	build.Dir = dir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("cmd/worker build failed: %v\n%s", err, out)
	}
}
