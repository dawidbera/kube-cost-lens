// Package main is the entry point for the kube-cost-lens binary.
package main

import (
	"os"

	"github.com/dawidbera/kube-cost-lens/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
