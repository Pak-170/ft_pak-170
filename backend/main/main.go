package main

import (
	"go_pakito/dbp"
	"go_pakito/web"
//	"net/http"
)

func main() {
	dbconnection, sqldb := dbp.OpenDatabaseConnection()
	defer sqldb.Close()

	var gormdb dbp.GormDB = dbp.GormDB{DB:dbconnection}
	web.Web(&gormdb)
}
