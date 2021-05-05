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
			fmt.Println("没有找到MAVEN的安装目录.如果不想设置 可以执行命令的时候 传入maven的目录参数(推荐此参数)!!!")
			pause()
		}
		//从配置文件中解析出 maven的本地仓库目录
		repoDir = parseConfig4repoDir(home)
		//如果从配置文件 还没有获取信息  那么默认就是查找当前用户目录下 的.m2文件夹

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

func clean(dir string) {
	//创建索引文件
	mcIndex := newMCIndex()

	list4del := make([]string, 0, 256)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		//标识是否是否可以被记录索引文件中
		flag := true

		if strings.Contains(path, "unknown") || strings.Contains(path, (string(os.PathSeparator))+"error") {
			fmt.Println("正在收集待删除文件/文件夹:", path)
			list4del = append(list4del, path)
			flag = false
		}
		//匹配删除 无用的 ${xxx-version}之类的无用文件夹
		if b, _ := regexp.MatchString("\\$\\{.*\\}", path); b {
			fmt.Println("正在收集待删除文件/文件夹:", path)
			list4del = append(list4del, path)
			flag = false
		}
		//删除文件后缀名为 lastUpdated的文件
		if strings.Contains(path, "lastUpdated") {
			fmt.Println("正在收集待删除文件/文件夹:", path)
			list4del = append(list4del, path)
			flag = false
		}

		if flag {
			//有资格参数 索引操作

			if mcIndex.CheckAndRecord(path, info) {

				//加入被删除的集合中
				fmt.Println("正在收集待删除文件/文件夹:", path)
				list4del = append(list4del, path)
			}
		}

		return err
	})

	if len(list4del) > 0 {
		fmt.Println("开始删除文件夹/文件!!!!")
		for _, x := range list4del {
			os.RemoveAll(x)
			fmt.Println("删除文件夹/文件:" + x)
		}
	}

	//删除无用的文件夹
	count := 0
	count, list4del = collectEmpty(dir)
	fmt.Printf("共有%d个文件\r\n", count)
	if len(list4del) > 0 {
		fmt.Println("开始删除文件夹/文件!!!!")
		for _, x := range list4del {
			os.RemoveAll(x)
			fmt.Println("删除文件夹/文件:" + x)
		}
	}

}

func GetMavenHome() string {
	home := os.Getenv("M2_HOME")
	if home == "" {
		home = os.Getenv("MAVEN_HOME")
	}
	return home
}
