package main

import (
	"go_pakito/dbp"
	"go_pakito/web"
//	"net/http"
)

func main() {
	dbconnection, sqldb := dbp.OpenDatabaseConnection()
	defer sqldb.Close()

	var gormdb web.GormDB = web.GormDB{DB:dbconnection}
	web.Web(&gormdb)
}
