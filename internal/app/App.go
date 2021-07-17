package app

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed "VERSION"
var version string

func Version() string {
	return version
}

func Exit(message string, code int) {
	fmt.Println(message)
	os.Exit(code)
}
