// file: db/init.go
package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite", "./power_reporting.db")
	if err != nil {
		log.Fatal(err)
	}

	// 测试数据库连接
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("数据库连接成功")
}

func GetDB() *sql.DB {
	return DB
}
