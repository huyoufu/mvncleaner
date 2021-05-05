package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/*var bytePool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 512)
		return &b
	},
}*/
func checkMRFileType(ori string) bool {
	filePath := filepath.Base(ori)
	idx := strings.LastIndexByte(filePath, '.')
	if idx > 0 {
		suffix := filePath[idx+1:]
		switch suffix {
		case "pom":
			return checkXmlFileType(ori)
		case "jar":
			return checkJarFileType(ori)
		default:
			fmt.Println("类型未知", ori)

		}

	}
	return true

}

func checkJarFileType(path string) bool {
	file, _ := os.Open(path)
	defer file.Close()

	buff := [4]byte{}

	file.Read(buff[:])
	bytesBuffer := bytes.NewBuffer(buff[:])

	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)

	i1 := int32(x)
	//文件头
	i2 := int32(0x04034b50)
	return i1 == i2

}

func checkXmlFileType(path string) bool {
	file, _ := os.Open(path)
	defer file.Close()

	reader := bufio.NewReader(file)

	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	} else {
		return strings.Contains(line, "<?xml") || strings.Contains(line, "<project>")
	}
}
