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

	if _, ok := os.LookupEnv("LIQUIBASE_HOME"); ok {
		liquibasepath = os.Getenv("LIQUIBASE_HOME")
	} else {
		// Find Liquibase Command
		out, err := exec.Command("which", "liquibase").CombinedOutput()
		if err != nil {
			fmt.Println("Unable to locate Liquibase.")
			os.Exit(1)
		}

		// Determine if Command is Symlink
		loc := strings.TrimRight(string(out), "\n")
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
	}

	commands.Execute(liquibasepath + "lib/")
}