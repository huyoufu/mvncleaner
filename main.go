package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {

	defaultTimerTrack.Start()
	dir := defaultFinder.GetRepo()
	clean(dir)
	defaultTimerTrack.End()
	printMetric()
	pause()
}
func printMetric() {
	metric := "****************************************************************************\n" +
		"**                                                                        **\n" +
		"**        该次清理统计信息如下:                                           **\n" +
		"**                                                                        **\n" +
		"**                                                                        **\n" +
		/*		"**        所有文件个数:" + fmt.Sprintf("%9d", defaultMCMetric.sumFileNum) + "                                            **\n" +
				"**        普通文件个数:" + fmt.Sprintf("%9d", defaultMCMetric.commonFileNum) + "                                            **\n" +*/
		"**        错误文件个数:" + fmt.Sprintf("%9d", defaultMCMetric.errFileNum) + "                                          **\n" +
		/*		"**        索引文件个数:" + fmt.Sprintf("%9d", defaultMCMetric.indexFileNum) + "                                            **\n" +
				"**        jar文件个数:" + fmt.Sprintf("%9d", defaultMCMetric.jarFileNum) + "                                            **\n" +*/
		"**        共花费了(ms):" + fmt.Sprintf("%9d", defaultTimerTrack.Cost()) + "                                          **\n" +
		"**                                                                        **\n" +
		"**                                                                        **\n" +
		"**                                                                        **\n" +
		"**        github: https://www.github.com/huyoufu/mvncleaner_mac              **\n" +
		"**        点个赞呗                                                        **\n" +
		"**                                                                        **\n" +
		"****************************************************************************\n"
	fmt.Println("")
	fmt.Println("")
	fmt.Println(metric)

}

func clean(dir string) {
	//创建索引文件
	mcIndex := newMCIndex()

	list4del := make([]string, 0, 256)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}

		//标识是否是否可以被记录索引文件中
		flag := true
		//if strings.Contains(path, "unknown") || strings.Contains(path, (string(os.PathSeparator))+"error") {
		if strings.Contains(path, "unknown") {
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
			//有资格参与 索引操作

			if mcIndex.CheckAndRecord(path, info) {

				//加入被删除的集合中
				fmt.Println("正在收集待删除文件/文件夹:", path)
				list4del = append(list4del, path)
			}
		}

		return err
	})

	//关闭索引文件
	mcIndex.destroy()

	if len(list4del) > 0 {
		fmt.Println("开始删除文件夹/文件!!!!")
		for _, x := range list4del {

			os.RemoveAll(x)
			defaultMCMetric.IncrErrFileNum()
			fmt.Println("删除文件夹/文件:" + x)
		}
	}

	//删除无用的文件夹

	_, list4del = collectEmpty(dir, dir)
	//fmt.Printf("共有%d个文件\r\n", count)
	if len(list4del) > 0 {
		fmt.Println("开始删除文件夹/文件!!!!")
		for _, x := range list4del {
			if dir == x {
				continue
			}
			os.RemoveAll(x)
			defaultMCMetric.IncrErrFileNum()
			fmt.Println("删除文件夹/文件:" + x)
		}
	}

}
