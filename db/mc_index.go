package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os/user"
	"path"
)

const DbName string = ".mvncleaner/mc_index.db"

var engine *sql.DB

func init() {
	var err error
	engine, err = sql.Open("sqlite3", getDBFilePath())
	if err != nil {
		fmt.Println(err)
	}
}
func Home() string {
	current, _ := user.Current()
	return current.HomeDir
}
func getDBFilePath() string {
	return path.Join(Home(), DbName)
}

func checkTableInit() {
	rows, _ := engine.Query("select count(*) as count from sqlite_master where name=?", "mc_index")
	defer rows.Close()
	if rows.Next() {
		var count int64
		rows.Scan(&count)
		fmt.Println(count)
	}
}
