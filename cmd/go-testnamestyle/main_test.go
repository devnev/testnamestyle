package main_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestNoFlags(t *testing.T) {
	binPath := compileBin(t)
	cmd := exec.Command(binPath, "./...")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Dir = "testdata"
	if err := cmd.Run(); err != nil {
		t.Fatalf("Expected linter to run and exit 0, got %v", err)
	}
}

func TestMatchAnythingRule(t *testing.T) {
	binPath := compileBin(t)
	cmd := exec.Command(binPath, "-rule", "Test.*", "./...")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Dir = "testdata"
	if err := cmd.Run(); err != nil {
		t.Fatalf("Expected linter to run and exit 0, got %v", err)
	}
}

func TestMatchOnlyTest(t *testing.T) {
	binPath := compileBin(t)
	cmd := exec.Command(binPath, "-rule", "Test", "./...")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Dir = "testdata"
	if err, exitErr := cmd.Run(), (*exec.ExitError)(nil); err == nil || !errors.As(err, &exitErr) || exitErr.ExitCode() != 3 {
		t.Fatalf("Expected linter to run and exit 3, got %v", err)
	}
}

func TestMatchOneGroup(t *testing.T) {
	binPath := compileBin(t)
	cmd := exec.Command(binPath, "-rule", "Test(.+)", "./...")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Dir = "testdata"
	if err := cmd.Run(); err != nil {
		t.Fatalf("Expected linter to run and exit 0, got %v", err)
	}
}

func compileBin(t *testing.T) string {
	binPath := filepath.Join(t.TempDir(), "go-testnamestyle")
	if err := exec.Command("go", "build", "-o", binPath, ".").Run(); err != nil {
		t.Fatalf("Failed to build binary for test: %v", err)
	}
	return binPath
}
