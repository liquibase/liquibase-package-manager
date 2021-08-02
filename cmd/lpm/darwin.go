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

	var liquibasehome string

	if _, ok := os.LookupEnv("LIQUIBASE_HOME"); ok {
		liquibasehome = os.Getenv("LIQUIBASE_HOME")
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
			link, err := filepath.EvalSymlinks(loc)
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
	//handles homebrew installation of liquibase
	if strings.Contains(liquibasehome, "Cellar") {
		liquibasehome = strings.Replace(liquibasehome, "/bin", "/libexec", 1)
	}
	if !strings.HasSuffix(liquibasehome, "/") {
		liquibasehome = liquibasehome + "/"
	}
	commands.Execute(liquibasehome, "/")
}