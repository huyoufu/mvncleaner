package main

import (
	"fmt"
	"testing"
)

func TestCollectEmpty(t *testing.T) {

	dir := "D:\\software\\maven-3.2.1\\repository"

	_, list4edl := collectEmpty(dir, dir)

	for _, filename := range list4edl {
		fmt.Println(filename)

	}
}
