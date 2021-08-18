// +build windows

package lpm

import (
	"strings"
)

func TrimCommandOutput(cmdout string) string {
	return strings.TrimRight(strings.Split(string(cmdout), "\n")[0], "\r")
}
