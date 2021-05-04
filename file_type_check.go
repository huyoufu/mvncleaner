package main

import "os"

func checkFileType(path string) bool {
	file, _ := os.Open(path)
	buff := make([]byte, 512)
	file.Read(buff)

	return true

}
