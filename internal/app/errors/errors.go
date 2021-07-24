package errors

import (
	"fmt"
	"os"
)

//Exit graceful exit with message
func Exit(message string, code int) {
	fmt.Println(message)
	os.Exit(code)
}

