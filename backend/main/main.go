package main

import (
	"go_pakito/dbp"
	"go_pakito/web"
)

func main() {
	dbconnection, sqldb := dbp.OpenDatabaseConnection()
	defer sqldb.Close()

	web.Web(dbconnection)
}
