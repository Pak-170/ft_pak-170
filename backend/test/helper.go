package test

import (
	"context"
	"go_pakito/dbp"
	"go_pakito/server"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
)

func setupTest(t *testing.T) (context.Context, *gin.Engine, *httptest.ResponseRecorder){
	ctx := context.Background()
	dbconnection, sqldb := dbp.OpenDatabaseConnection()

	tx := dbconnection.Begin()

	t.Cleanup(func(){
		tx.Rollback()
		sqldb.Close()
	})
	r := server.NewRouter(ctx, tx)
	w := httptest.NewRecorder()

	return ctx, r, w
}