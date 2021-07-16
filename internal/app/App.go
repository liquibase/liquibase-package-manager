package app

import (
	"package-manager/internal/app/commands"
)

func Exec(cp string) {
	commands.Execute(cp)
}