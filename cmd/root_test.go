package cmd

import (
	"bytes"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"version"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing version command: %v", err)
	}
}

func TestRootFlags(t *testing.T) {
	if rootCmd.Flag("demo") == nil {
		t.Error("expected --demo flag on rootCmd")
	}
	if rootCmd.Flag("namespace") == nil {
		t.Error("expected --namespace flag on rootCmd")
	}
	if rootCmd.Flag("rates") == nil {
		t.Error("expected --rates flag on rootCmd")
	}
	if rootCmd.Flag("refresh-interval") == nil {
		t.Error("expected --refresh-interval flag on rootCmd")
	}
}
