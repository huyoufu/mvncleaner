package main

import (
	"encoding/xml"
	"fmt"
	"github.com/huyoufu/mvncleaner/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {

	dir := getRepoDir()
	clean(dir)

}
func getRepoDir() string {
	repoDir := ""
	args := os.Args
	if len(args) > 1 {
		repoDir = args[1]
	} else {

		//没有指明路径参数 获取maven的安装目录
		home := GetMavenHome()
		if home == "" {
			fmt.Println("没有找到MAVEN的安装目录.如果不想设置 可以执行命名的时候 传入maven的目录参数(推荐此参数)!!!")
			os.Exit(0)
		}

		repoDir = parseConfig4repoDir(home)
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

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if strings.Contains(path, "unknown") || strings.Contains(path, "/error") {
			fmt.Println("正在删除文件/文件夹:", path)
			delErr := os.RemoveAll(path)
			if delErr != nil {
				fmt.Println(delErr)
			}
		}
		if b, _ := regexp.MatchString("\\$\\{.*\\}", path); b {
			fmt.Println("正在删除文件/文件夹:", path)
			os.RemoveAll(path)
		}

		if strings.Contains(path, "lastUpdated") {
			fmt.Println("正在删除文件/文件夹:", path)
			delErr := os.Remove(path)
			if delErr != nil {
				fmt.Println(delErr)
			}
		}

		return err
	})

}
func GetMavenHome() string {
	home := os.Getenv("M2_HOME")
	if home == "" {
		home = os.Getenv("MAVEN_HOME")
	}
	return home
}
