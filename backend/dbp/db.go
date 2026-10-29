package dbp

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"database/sql"
	"go_pakito/models"
)

const DB_TIMEOUT = 5


// SELECT ---------------------------------------------------------------------

func SelectUsersOfChat(ctx context.Context, db *gorm.DB, chatID uint) ([]models.User, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, DB_TIMEOUT * time.Second)
	defer cancel()
	
	usersInChat, err := gorm.G[models.User](db).
		Raw("SELECT u.id, u.username, u.alias FROM users AS u JOIN chat_user as c ON c.user = u.id WHERE c.chat = ?", chatID).
		Find(ctxWithTimeout)
	return usersInChat, err
}

func SelectAllChatsOfUser(ctx context.Context, db *gorm.DB, userID uint) ([]models.Chat, []models.ChatView, error) { //! Optimizar las consultas
	ctxWithTimeout, cancel := context.WithTimeout(ctx, DB_TIMEOUT * time.Second)
	defer cancel()

	chatUserRelation, err := gorm.G[models.ChatUserRelation](db).Where("user = ?", userID).Find(ctx)

	if err != nil {
		return nil, nil, err
	}

	var chatsMsgs	[]models.ChatView = make([]models.ChatView, 0)
	var chats		[]models.Chat = make([]models.Chat, 0)
	for _, chatUserIsIn := range chatUserRelation {
		//1. Get the chat info:
		chatInfo, err := gorm.G[models.Chat](db).Where("id = ?", chatUserIsIn.ChatID).First(ctx)
		if err != nil {
			return nil, nil, err
		}

		//2. Get the list of users inside said chat:
		usersInChat, err := SelectUsersOfChat(ctx, db, chatUserIsIn.ChatID)
		if err != nil {
			return nil, nil, err
		}
		chatInfo.Participants = usersInChat

		//3. Append the chatInfo we got to the chats slice we will send to the frontend
		chats = append(chats, chatInfo)

		//4. Find the msgs of that chat to send to the frontend as prior messages.
		oneChatMsgs, err := gorm.G[models.ChatView](db).Where("chat_id = ?", chatUserIsIn.ChatID).Find(ctxWithTimeout)
		if err != nil {
			return nil, nil, err
		}
		chatsMsgs = append(chatsMsgs, oneChatMsgs...)
	}
	return chats, chatsMsgs, err
}

// func	SelectAnyWhere[T any](ctx context.Context, db *gorm.DB, where string) ([]T, error) {
// 	ctxWithTimeout, cancel := context.WithTimeout(ctx, DB_TIMEOUT * time.Second)
// 	defer cancel()

// 	tableRegisters, err := gorm.G[T](db).Where(where).Find(ctxWithTimeout)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return tableRegisters, err
// }


func	SelectAny[T any](ctx context.Context, db *gorm.DB) ([]T, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, DB_TIMEOUT * time.Second)
	defer cancel()

	tableRegisters, err := gorm.G[T](db).Find(ctxWithTimeout)

	if err != nil {
		return nil, err
	}
	return tableRegisters, err
}

func	SelectChatMsgs(ctx context.Context, db *gorm.DB, chatID uint) ([]models.ChatView, error) {


	ctxWithTimeout, cancel := context.WithTimeout(ctx, DB_TIMEOUT * time.Second)
	defer cancel()
	var err error

	chats, err := gorm.G[models.ChatView](db).Where("chat_id = ?", chatID).Find(ctxWithTimeout)

	return chats, err
}

// INSERT ---------------------------------------------------------------------

func	InsertMsg(ctx context.Context, db *gorm.DB, msgInChat *models.ChatView) error {
	var	msgInDB models.Message = models.Message{
		ChatID: msgInChat.ChatID,
		Msg: msgInChat.Msg,
		Time: msgInChat.SendTime,
	}
	ctxWithTimeout, cancel := context.WithTimeout(ctx, DB_TIMEOUT * time.Second)
	defer cancel()
	var	err error

	err = gorm.G[models.User](db).Select("id").Where("username = ?", msgInChat.Username).Scan(ctxWithTimeout, &msgInDB.SenderID)
	if err != nil {
		log.Printf("ERROR CREATING THE MSG MODEL: %v\n", err)
	}
	err = gorm.G[models.Message](db).Create(ctx, &msgInDB)
	if err != nil {
		log.Printf("ERROR INSERTING INTO DB: %v\n", err)
	}
	return err
}

func InsertUserInChat(ctx context.Context, db *gorm.DB, chat uint, user uint) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, DB_TIMEOUT * time.Second)
	defer cancel()
	var userInChat models.ChatUserRelation = models.ChatUserRelation{
		ChatID: chat,
		UserID: user,
	}
	return gorm.G[models.ChatUserRelation](db).Create(ctxWithTimeout, &userInChat)
}

func InsertChat(ctx context.Context, db *gorm.DB, chatName string, owner uint) (uint, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, DB_TIMEOUT * time.Second)
	defer cancel()
	var chat models.Chat = models.Chat{
		ID: 0,
		Name: chatName,
		Owner: owner,
	}
	err := gorm.G[models.Chat](db).Create(ctxWithTimeout, &chat)
	return chat.ID, err
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