package main

import (
	"encoding/xml"
	"fmt"
	"github.com/huyoufu/mvncleaner/config"
	"github.com/huyoufu/mvncleaner/log"
	"io/ioutil"
	"os"
	"path/filepath"
)

type MCRFinder struct {
	_index    map[string]int64
	IsNew     bool
	Size      int64
	indexFile *os.File
}

var defaultFinder MCRFinder = MCRFinder{}

//获取仓库的文件夹名字
func (f *MCRFinder) GetRepo() string {
	var repoDir string = ""
	args := os.Args
	if len(args) > 1 {
		repoDir = args[1]
	} else {

		//没有指明路径参数 获取maven的安装目录
		home := GetMavenHome()
		log.Info("根据环境变量MAVEN_HOME/M2_HOME 得到maven软件的家目录:" + home)
		if home == "" {
			fmt.Println("没有找到MAVEN的安装目录.如果不想设置MAVEN_HOME 可以执行命令的时候 传入MAVEN的本地仓库目录!!!")
			pause()
			/*reader := bufio.NewReader(os.Stdin)
			fmt.Print("请输入maven的本地仓库目录:> ")

			repoDir,_= reader.ReadString('\n')
			// convert CRLF to LF
			repoDir = strings.Replace(repoDir, "\n", "", -1)

			os.Stdin.Close()*/

		} else {
			//从配置文件中解析出 maven的本地仓库目录
			repoDir = parseConfig4repoDir(home)
			//如果从配置文件 还没有获取信息  那么默认就是查找当前用户目录下 的.m2文件夹
		}
	}

	_, err := os.Lstat(repoDir)
	if err != nil {
		fmt.Println("找不到该文件目录:" + repoDir)
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

//获取环境变量中的MAVEN_HOME
func GetMavenHome() string {
	home := os.Getenv("MAVEN_HOME")
	if home == "" {
		home = os.Getenv("M2_HOME")
	}
	return home
}
