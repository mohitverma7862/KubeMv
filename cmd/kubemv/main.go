package main

import (
	"fmt"
	"os"
)

const version = "0.0.0-dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "--version") {
		fmt.Printf("kubeMV %s\n", version)
		return
	}
	fmt.Println("kubeMV — CLI-based Kubernetes tool")
	fmt.Println("Usage: kubemv version")
}
