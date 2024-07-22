package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/liquibase/liquibase-package-manager/internal/app/commands"
)

func main() {

	var liquibasehome string

	if _, ok := os.LookupEnv("LIQUIBASE_HOME"); ok {
		liquibasehome = os.Getenv("LIQUIBASE_HOME")
	} else {
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
			liquibasehome, _ = filepath.Split(link)
		} else {
			// Not Symlink
			liquibasehome, _ = filepath.Split(loc)
		}
	}
	if !strings.HasSuffix(liquibasehome, "\\") {
		liquibasehome = liquibasehome + "\\"
	}
	commands.Execute(liquibasehome, "\\")
}
