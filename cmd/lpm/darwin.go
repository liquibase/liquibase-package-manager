package main

import (
	"fmt"
	"os"
	"package-manager/internal/app"
)

func main() {
	jh := os.Getenv("JAVA_HOME")
	cp := os.Getenv("CLASSPATH")

	if jh == "" {
		fmt.Println("JAVA_HOME not found.")
		os.Exit(1)
	}
	if cp == "" {
		cp = jh + "/lib/"
	}

	app.Exec(cp)
}