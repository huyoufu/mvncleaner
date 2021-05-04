package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const INDEXFILENAME = "index.mc"

//创建索引文件
func getIndexFile() (*os.File, bool) {
	userHome, _ := os.UserHomeDir()
	indexPath := path.Join(userHome, "mvncleaner")

	_, err := os.Stat(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			//创建索引目录
			err := os.Mkdir(indexPath, os.ModePerm)
			if err != nil {
				fmt.Printf("无法创建索引目录:%s\n", indexPath)
			}
		}
	}
	indexFilePath := path.Join(indexPath, INDEXFILENAME)

	var indexFile *os.File
	_, err2 := os.Stat(indexFilePath)
	if err2 != nil {
		if os.IsNotExist(err2) {
			//创建索引目录
			indexFile, err2 = os.Create(indexFilePath)
			if err2 != nil {
				fmt.Printf("无法创建索引文件:%s\n", indexFilePath)
			}
			writeIndexInitInfo(indexFile)
			return indexFile, true
		}
	} else {
		indexFile, _ = os.OpenFile(indexFilePath, os.O_RDWR, 0666)
	}
	return indexFile, false
}

//读取索引
func readIndex() map[string]int64 {
	maps := make(map[string]int64, 256)

	file, isNew := getIndexFile()
	if isNew {
		return maps
	}

	reader := bufio.NewReader(file)

	num := 13
	for {
		line, err := reader.ReadString('\n')
		if num > 0 {
			num--
			continue
		}

		if err == io.EOF {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		words := strings.Split(line, "=")
		millsSecond, _ := strconv.ParseInt(words[1], 10, 64)

		maps[words[0]] = millsSecond
	}
	return maps
}
func readIndex1() map[string]time.Time {
	maps := make(map[string]time.Time, 256)

	file, isNew := getIndexFile()
	if isNew {
		return maps
	}
	reader := bufio.NewReader(file)

	num := 13
	for {
		line, err := reader.ReadString('\n')
		if num > 0 {
			num--
			continue
		}

		if err == io.EOF {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		words := strings.Split(line, "=")
		millsSecond, _ := strconv.ParseInt(words[1], 10, 64)

		maps[words[0]] = time.Unix(millsSecond, 0)
	}
	return maps
}

//生成索引文件 初始化信息
func writeIndexInitInfo(indexFile *os.File) {

	now := time.Now().Format("2006-01-02 15:04:05")
	info := "****************************************************************************\n" +
		"**                                                                        **\n" +
		"**                                                                        **\n" +
		"**        this file is create by mvcleaner                                **\n" +
		"**        do not modify this file                                         **\n" +
		"**        create file time: " + now + "                           **\n" +
		"**        the program author is huyoufu                                   **\n" +
		"**        mailto:371778981@qq.com                                         **\n" +
		"**        websit: http://www.jk1123.com                                   **\n" +
		"**        github: https://www.github.com/huyoufu/mvcleaner                **\n" +
		"**                                                                        **\n" +
		"**                                                                        **\n" +
		"****************************************************************************\n"

	indexFile.WriteString(info)

}
