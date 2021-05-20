package main

import (
	"fmt"
	"testing"
)

func TestFindRepo(t *testing.T) {
	repo := defaultFinder.GetRepo()
	fmt.Println(repo)

}
