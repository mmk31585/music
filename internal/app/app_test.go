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

	// Verify that the API and worker entry points build without conflict.
	// Build into a temp dir: a multi-package `go build` writes an executable per
	// main package, which would otherwise overwrite the tracked ./api and ./worker
	// binaries in the repo root.
	build := exec.Command("go", "build", "-o", t.TempDir(), "./cmd/api", "./cmd/worker")
	build.Dir = dir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("cmd build failed: %v\n%s", err, out)
	}
}
