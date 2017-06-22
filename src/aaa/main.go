package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "E:\\dnapersondb\\src\\database\\person23.db")
	checkErr(err)

	//插入数据
	_, err = db.Exec("insert into sys_dict(id) values ('1')")
	checkErr(err)
	db.Close()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
