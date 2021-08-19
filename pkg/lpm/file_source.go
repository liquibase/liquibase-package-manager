package lpm

type FileSource byte

const GlobalFiles FileSource = 'g'
const LocalFiles FileSource = 'l'

var fileSources = map[bool]FileSource{
	true:  GlobalFiles,
	false: LocalFiles,
}

func GetFileSource() (fs FileSource) {
	return fileSources[cliArgs.Global]
}
