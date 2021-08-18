package main

import (
	_ "embed" // Embed Import for Package Files
	"os"
	"package-manager/pkg/lpm"
)

func main() {

	err := lpm.Execute("/")
	if err != nil {
		lpm.ShowUserError(err)
		os.Exit(1)
	}
}
