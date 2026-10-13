package main

import (
	"go_pakito/dbp"
	"go_pakito/web"
//	"net/http"
	"context"
	// "github.com/markbates/goth"
	"golang.org/x/oauth2"
)


func initOauth() {
	UID := "u-s4t2ud-ec20f93efec6432580f6f25ce5afee4f9d73bdacfc665d18709f1bc6407ab786"
	SECRET := "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"

	var ctx context.Context = context.Background()
	var conf *oauth2.Config = &oauth2.Config{
		ClientID: UID,
		ClientSecret: SECRET}


}

func main() {
	dbconnection, sqldb := dbp.OpenDatabaseConnection()
	defer sqldb.Close()

	var gormdb dbp.GormDB = dbp.GormDB{DB:dbconnection}
	web.Web(&gormdb)
}
