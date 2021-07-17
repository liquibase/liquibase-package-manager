package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"package-manager/internal/app/commands"
	"path/filepath"
	"strings"
)

func main() {

	var liquibasepath string

	// Find Liquibase Command
	out, err := exec.Command("where", "liquibase").CombinedOutput()
	if err != nil {
		fmt.Println("Unable to locate Liquibase.")
		os.Exit(1)
	}

	// Determine if Command is Symlink
	loc := strings.TrimRight(strings.Split(string(out), "\n")[0], "\r")
	fi, err := os.Lstat(loc)
	if err != nil {
		log.Fatal(err)
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		link, err := os.Readlink(loc)
		if err != nil {
			log.Fatal(err)
		}
		// Is Symlink
		liquibasepath, _ = filepath.Split(link)
	} else {
		// Not Symlink
		liquibasepath, _ = filepath.Split(loc)
	}

	commands.Execute(liquibasepath + "lib\\")
}