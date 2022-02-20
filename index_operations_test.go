package main

import (
	"fmt"

	"testing"
)

func TestGetIndexFile(t *testing.T) {

	file, _ := getIndexFile()
	fmt.Println(file.Name())

}
