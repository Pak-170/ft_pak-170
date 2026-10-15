package main

import (
	"go_pakito/web"
	"go_pakito/dbp"
)
func main() {
	dbconnection, sqldb := dbp.OpenDatabaseConnection()
	defer sqldb.Close()

	var gormdb dbp.GormDB = dbp.GormDB{DB:dbconnection}
	web.Web(&gormdb)
}
