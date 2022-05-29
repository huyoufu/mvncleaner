package db

import (
	"fmt"
	"testing"
)

func TestMustCreateIndex(t *testing.T) {
	rows, _ := engine.Query("select count(*) as count from sqlite_master where name=?", "mc_index")
	defer rows.Close()
	if rows.Next() {
		var count int64
		rows.Scan(&count)
		fmt.Println(count)
	}

}
