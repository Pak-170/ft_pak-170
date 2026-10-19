package dbp

import (
	"context"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"database/sql"
	"go_pakito/models"
)

const ALL = -1

// SELECT ---------------------------------------------------------------------

/* func SelectUserByUsername(db *gorm.DB, username string) (models.User, error) {
	ctx := context.Background()
	return gorm.G[models.User](db).Where("username = ?", username).First(ctx)
}

func SelectChatByChatname(db *gorm.DB, chatname string) ([]models.ChatView, error) {
	ctx := context.Background()
	fmt.Println(gorm.G[models.ChatView](db).First(ctx))
	return gorm.G[models.ChatView](db).Where("name = ?", chatname).Find(ctx)
} */

func SelectAllChatsOfUser(db *gorm.DB, userID uint) ([]models.ChatView, error) { //! Optimizar las consultas
	ctx := context.Background() //! investigar mejor esto de context
	chatUserRelation, err := gorm.G[models.ChatUserRelation](db).Where("user = ?", userID).Find(ctx)

	if err != nil {
		return nil, err
	}

	var chatsMsgs []models.ChatView
	for _, chatUserIsIn := range chatUserRelation {
		oneChatMsgs, err := gorm.G[models.ChatView](db).Where("chat_id = ?", chatUserIsIn.ChatID).Find(ctx)
		if err != nil {
			return nil, err
		}
		chatsMsgs = append(chatsMsgs, oneChatMsgs...)
	}
	return chatsMsgs, err
}

func	SelectChatIds(db *gorm.DB, userID int) ([]int, error) {
	ctx := context.Background()

	var	chatIds []int
	var err error

	if userID == ALL {
		err = gorm.G[models.Chat](db).Select("id").Scan(ctx, &chatIds)
	} else {
		err = gorm.G[models.ChatUserRelation](db).Select("chat").Where("user = ?", userID).Scan(ctx, &chatIds)
	}
	return chatIds, err
}

// INSERT ---------------------------------------------------------------------

func	InsertMsg(db *gorm.DB, msgInChat *models.ChatView) error {
	ctx := context.Background()
	var	msgInDB models.Message = models.Message{
		ChatID: msgInChat.ChatID,
		Msg: msgInChat.Msg,
		Time: msgInChat.SendTime,
	}
	var	err error
	err = gorm.G[models.User](db).Select("id").Where("username = ?", msgInChat.Username).Scan(ctx, &msgInDB.SenderID)
	if err != nil {
		log.Printf("ERROR CREATING THE MSG MODEL: %v\n", err)
	}
	err = gorm.G[models.Message](db).Create(ctx, &msgInDB)
	if err != nil {
		log.Printf("ERROR INSERTING INTO DB: %v\n", err)
	}
	return err
}

// ----------------------------------------------------------------------------

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
			Logger: logger.Default.LogMode(logger.Info),
		})
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