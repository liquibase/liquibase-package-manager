package main

import (
	_ "embed" // Embed Import for Package Files
	"os"
	"package-manager/pkg/lpm"
	"package-manager/pkg/lpm/cmd"
)

func main() {

	err := cmd.Execute("/")
	if err != nil {
		lpm.ShowUserError(err)
		os.Exit(1)
	}
}
