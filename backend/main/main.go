package main

import (
	"go_pakito/dbp"
	"go_pakito/server"
	"go_pakito/validation"
	"context"
)

func main() {
	ctx := context.Background()
	dbconnection, sqldb := dbp.OpenDatabaseConnection()
	defer sqldb.Close()
	
	validation.InitValidator()
	server.InitServer(ctx, dbconnection)
}
