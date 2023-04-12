package file

var fileName = ""
var mainFileName = ""

func SetMainFileName(name string) {
	mainFileName = name
}

func GetMainFileName() string {
	return mainFileName
}

func SetFileName(name string) {
	fileName = name
}

func GetFileName() string {
	return fileName
}
