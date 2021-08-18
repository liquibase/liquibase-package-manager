// +build darwin

package lpm

import (
	"strings"
)

//goland:noinspection GoUnusedExportedFunction
func TrimCommandOutput(cmdout string) string {
	return strings.TrimRight(string(cmdout), "\n")
}
