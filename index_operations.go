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

const IndexFileName = "index.mc"

type MCIndex struct {
	_index    map[string]int64
	IsNew     bool
	Size      int64
	indexFile *os.File
}

func newMCIndex() *MCIndex {
	mcIndex := new(MCIndex)
	indexFile, isNew := getIndexFile()
	mcIndex.IsNew = isNew
	mcIndex.indexFile = indexFile

	mcIndex.readIndex(indexFile)
	//seek读写位置到
	if !isNew {
		info, _ := mcIndex.indexFile.Stat()
		mcIndex.Size = info.Size()
		indexFile.Seek(info.Size(), io.SeekStart)
	}

	return mcIndex
}

//检查该文件是否在索引中 是否过期 类型是否正确
//如果返回false表示无需删除
//如果返回true标识要删除
func (i *MCIndex) CheckAndRecord(filePath string, info os.FileInfo) bool {
	if info.IsDir() {
		//目录直接略过
		return false
	}

	if i.IsNew {
		//第一次新建的索引文件 所以直接判断 文件类型是否正确
		//如果文件类型正确 添加索引
		fileTypeIsCorrect := checkMRFileType(filePath)
		if fileTypeIsCorrect {
			i.writeIndex(filePath, info)
			return false
		} else {
			return true
		}
	}

	//如果是文件  则查看索引中是否存在!!
	modify := i.checkModify(filePath, info)
	if modify {
		//检查文件内容
		fileTypeIsCorrect := checkMRFileType(filePath)
		if !fileTypeIsCorrect {
			return true
		} else {
			i.writeIndex(filePath, info)
			return false
		}

	}
	return false
}
func (i *MCIndex) checkModify(path string, info os.FileInfo) bool {

	//有索引文件就得判断了
	time := i._index[path]
	if time != 0 {
		//有值的话 就判断更新时间
		if info.ModTime().Unix() > time {
			//说明有修改了
			return true
		}
	}
	return false
}

//销毁索引对象
func (i *MCIndex) destroy() {
	i.indexFile.Close()
	i._index = nil
}

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
	indexFilePath := path.Join(indexPath, IndexFileName)

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
func (i *MCIndex) readIndex(file *os.File) {
	i._index = make(map[string]int64, 256)

	if i.IsNew {
		return
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
		i._index[words[0]] = millsSecond
	}

}

func (i *MCIndex) writeIndex(filePath string, info os.FileInfo) {

	time := info.ModTime().Unix()

	_, err := i.indexFile.WriteString(filePath + "=" + strconv.FormatInt(time, 10) + "\n")
	if err != nil {
		fmt.Println("写入内容失败了!", filePath, err)
	}
}

//回头再说
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
		"**        website: http://www.jk1123.com                                  **\n" +
		"**        github: https://www.github.com/huyoufu/mvcleaner                **\n" +
		"**                                                                        **\n" +
		"**                                                                        **\n" +
		"****************************************************************************\n"

	indexFile.WriteString(info)

}
