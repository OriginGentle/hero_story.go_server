package base

import (
	"database/sql"
	"time"
)

import _ "github.com/go-sql-driver/mysql"

var MysqlDB *sql.DB

const DriverName = "mysql"
const DataSourceName = "root:Ycb@990121@tcp(139.155.0.174:3236)/hero_story"

func init() {
	var mysqlErr error

	MysqlDB, mysqlErr = sql.Open(DriverName, DataSourceName)

	if nil != mysqlErr {
		panic(mysqlErr)
	}

	MysqlDB.SetMaxOpenConns(128)
	MysqlDB.SetMaxIdleConns(16)
	MysqlDB.SetConnMaxIdleTime(2 * time.Minute)

	if mysqlErr = MysqlDB.Ping(); nil != mysqlErr {
		panic(mysqlErr)
	}
}
