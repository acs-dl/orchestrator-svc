package main

import (
	"os"

	"gitlab.com/distributed_lab/acs/orchestrator/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
