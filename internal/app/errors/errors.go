package errors

import (
	"fmt"
	"os"
)

func Exit(message string, code int) {
	fmt.Println(message)
	os.Exit(code)
}

