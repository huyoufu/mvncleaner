package main

import (
	"fmt"

	"testing"
)

func TestGetIndexFile(t *testing.T) {

	file := getIndexFile()
	fmt.Println(file)

}
func TestReadIndex(t *testing.T) {

	index := readIndex()
	fmt.Println(len(index))
	//遍历
	for key, value := range index {
		fmt.Printf("key: %s -----------------value: %d\n", key, value)
	}
}
