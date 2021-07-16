package app

import (
	"package-manager/internal/app/commands"
)

func Exec(cp string) {
	//files, err := ioutil.ReadDir(cp)
	//if err != nil {
	//	panic(err)
	//}
	//for _, f := range files {
	//	fmt.Println(f.Name())
	//}
	commands.Execute()
}