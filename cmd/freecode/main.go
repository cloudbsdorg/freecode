package main

import (
	"fmt"
	"os"

	"github.com/freecode/freecode/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
