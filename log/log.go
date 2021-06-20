package log

import "fmt"

type logger struct {
}

var log logger

func init() {
	fmt.Println("初始化")
}
func Info(info interface{}) {
	fmt.Println(info)
}
