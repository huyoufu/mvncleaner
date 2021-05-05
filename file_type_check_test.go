package main

import (
	"fmt"
	"testing"
)

func TestCheckFileType(t *testing.T) {

	fileType := checkMRFileType("D:\\software\\maven-3.2.1\\repository" +
		"\\com\\alibaba\\csp\\sentinel-core\\1.8.1\\sentinel-core-1.8.1.pom")

	fmt.Println(fileType)

}
