package dbp

import (
	"context"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"database/sql"
	"go_pakito/models"
)

type GormDB struct {
	DB *gorm.DB
}

// ------------------------ select by

func SelectUserByUsername(db *gorm.DB, username string) (models.User, error) {
	ctx := context.Background()
	return gorm.G[models.User](db).Where("username = ?", username).First(ctx)
}

func SelectChatByChatname(db *gorm.DB, chatname string) ([]models.ChatView, error) {
	ctx := context.Background()
	fmt.Println(gorm.G[models.ChatView](db).First(ctx))
	return gorm.G[models.ChatView](db).Where("name = ?", chatname).Find(ctx)
}

// ------------------------

func buildDSN() string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, name,
	)
}

func OpenDatabaseConnection() (*gorm.DB, *sql.DB) {
	dsn := buildDSN()

	db, err := gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			TranslateError: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			}})
	if err != nil {
		log.Fatalf("open db error: %v", err)
	}

	sqlDB, err := db.DB() //!ESTA VARIABLE GESTIONA EL POOL DE CONEXIONES, POR AHORA VOY A DEJAR LOS VALORES POR DEFECTO.
	if err != nil {
		log.Fatalf("get DB error: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	} else {
		fmt.Println("DB OK")
	}

	return db, sqlDB
}