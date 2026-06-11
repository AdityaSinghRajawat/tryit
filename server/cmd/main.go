// Package main is the tryit local server entry point. Thin: defers all
// initialisation logic to app.Run().
package main

import (
	"fmt"
	"os"
)

func main() {
	if err := Run(); err != nil {
		fmt.Fprintln(os.Stderr, "tryit: fatal:", err)
		os.Exit(1)
	}
}
