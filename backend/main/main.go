package main

import (
	"go_pakito/dbp"
	"go_pakito/server"
	"context"
)

func main() {
	ctx := context.Background()
	dbconnection, sqldb := dbp.OpenDatabaseConnection()
	defer sqldb.Close()

	server.InitServer(ctx, dbconnection)
}
