package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func collectEmpty(path string) (int, []string) {
	var count int = 0
	list4edl := make([]string, 0, 2)
	fis, _ := readDirNames(path)
	for _, f := range fis {
		//fmt.Println(f.Name())
		if f.IsDir() {
			ca, del4a := collectEmpty(filepath.Join(path, f.Name()))
			count += ca
			list4edl = append(list4edl, del4a...)
		} else {
			count++
		}
	}
	if count == 0 {
		fmt.Println("正在收集空文件夹:" + path)
		list4edl = append(list4edl, path)
	}
	return count, list4edl
}

func readDirNames(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	fis, err := f.Readdir(-1)

	f.Close()
	return fis, err
}
