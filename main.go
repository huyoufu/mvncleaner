package main

import (
	"encoding/xml"
	"fmt"
	"github.com/huyoufu/mvncleaner/config"
	"github.com/nsf/termbox-go"
	_ "github.com/nsf/termbox-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	termbox.SetCursor(0, 0)
	termbox.HideCursor()
}
func pause() {
	fmt.Println("请按任意键继续...")
Loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			break Loop
		}
	}
	//执行结束后关键 termbox
	termbox.Close()
}

type MavenRepoNoFindError struct {
}

func (e *MavenRepoNoFindError) Error() string {
	return "maven的仓库目录不存在!请告知maven仓库位置"
}

func main() {
	start := time.Now()
	dir := getRepoDir()
	clean(dir)
	end := time.Now()
	i := end.Unix() - start.Unix()
	fmt.Printf("清理完成!!!!!,共花费了%d秒\n", i)
	pause()

}
func getRepoDir() (repoDir string) {
	args := os.Args
	if len(args) > 1 {
		repoDir = args[1]
	} else {

		//没有指明路径参数 获取maven的安装目录
		home := GetMavenHome()
		if home == "" {
			fmt.Println("没有找到MAVEN的安装目录.如果不想设置 可以执行命名的时候 传入maven的目录参数(推荐此参数)!!!")
			pause()

		}

		repoDir = parseConfig4repoDir(home)

	}
	_, err := os.Lstat(repoDir)
	if err != nil {
		fmt.Printf("找不到该文件目录:%s\n", repoDir)
		pause()
	}
	return repoDir
}

func parseConfig4repoDir(home string) string {
	configFile := filepath.Join(home, "conf", "settings.xml")
	//接下xml文件
	f, err := os.Open(configFile)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	v := config.Settings{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	return v.LocalRepository
}
func add(ori []string, source string) []string {
	if len(ori) == 0 {
		ori = append(ori, source)
	} else {
		flag := false
		for _, x := range ori {
			if strings.Contains(source, x) {
				flag = true
				break
			}
		}
		if !flag {
			ori = append(ori, source)
		}
	}

	return ori
}

func clean(dir string) {
	//没被删除 就可以保存了
	indexFile, isNew := getIndexFile()
	defer indexFile.Close()
	index := readIndex()
	list4delx := make([]string, 0)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		//标识是否是否可以被记录索引文件中
		flag := true

		if strings.Contains(path, "unknown") || strings.Contains(path, (string(os.PathSeparator))+"error") {
			fmt.Println("正在收集待删除文件/文件夹:", path)
			list4delx = add(list4delx, path)
			flag = false
		}
		//匹配删除 无用的 ${xxx-version}之类的无用文件夹

		if b, _ := regexp.MatchString("\\$\\{.*\\}", path); b {
			fmt.Println("正在收集待删除文件/文件夹:", path)
			list4delx = add(list4delx, path)
			flag = false
		}

		if strings.Contains(path, "lastUpdated") {
			fmt.Println("正在收集待删除文件/文件夹:", path)
			list4delx = add(list4delx, path)
			flag = false
		}

		if !isNew {

			if checkModify(index, path, info) {
				//这就是修改了!!!!!
				fmt.Println("这个文件", path, "被修改了!!!")

			}
		}

		if flag {

			//如果是是文件 且新创建的索引文件
			if isNew && !info.IsDir() {
				//fmt.Println("当前文件名字",path)

				//判断文件类型是否正确

				//time := info.ModTime().Format("2006-01-02 15:04:05")
				time := info.ModTime().Unix()

				_, err := indexFile.WriteString(path + "=" + strconv.FormatInt(time, 10) + "\n")
				if err != nil {
					fmt.Println("写入内容失败了!", path, err)
				}
			} else {

			}
		}

		return err
	})

	if len(list4delx) > 0 {
		fmt.Println("开始删除文件夹/文件!!!!")
		for _, x := range list4delx {
			os.RemoveAll(x)
			fmt.Println("删除文件夹/文件:" + x)
		}
	}

	//删除无用的文件夹
	count := 0
	count, list4delx = collectEmpty(dir)
	fmt.Printf("共有%d个文件\r\n", count)
	if len(list4delx) > 0 {
		fmt.Println("开始删除文件夹/文件!!!!")
		for _, x := range list4delx {
			os.RemoveAll(x)
			fmt.Println("删除文件夹/文件:" + x)
		}
	}

}
func checkModify(index map[string]int64, path string, info os.FileInfo) bool {
	//有索引文件就得判断了
	time := index[path]
	if time != 0 {
		//有值的话 就判断更新时间
		if info.ModTime().Unix() > time {
			//说明有修改了
			return true
		}
	}
	return false
}

func GetMavenHome() string {
	home := os.Getenv("M2_HOME")
	if home == "" {
		home = os.Getenv("MAVEN_HOME")
	}
	return home
}
