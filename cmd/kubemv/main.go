package main

import (
	"fmt"
	"os"

	"github.com/mohitverma7862/KubeMv/internal/tui"
)

const version = "0.0.0-dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "--version") {
		fmt.Printf("kubemv %s\n", version)
		return
	}
	if err := tui.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
