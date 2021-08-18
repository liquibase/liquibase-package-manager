package lpm

type CliArgs struct {
	Category string
	Global   bool
}

var cliArgs CliArgs

func init() {
	cliArgs = CliArgs{}
}
