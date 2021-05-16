package umscraper

import (
	"io/ioutil"
	"os"
)

func ReadFile(filename string) []byte {
	file, err := os.Open(filename)
	checkError("Cannot open file", err)
	buf, err := ioutil.ReadAll(file)
	checkError("Cannot ReadAll", err)
	return buf
}

func WriteFile(content []byte, filename string) {
	err := ioutil.WriteFile(filename, content, 0644)
	checkError("Cannot write file", err)
}

func CreateDirectoryIfNotExist(dir string, fileMode os.FileMode) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, fileMode)
	}
}
