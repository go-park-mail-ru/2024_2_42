package repository

import (
	"database/sql"
	"fmt"

	"pinset/mailer-service/mailer"
	"pinset/mailer-service/models"
	"pinset/mailer-service/usecase"

	"github.com/sirupsen/logrus"
)

type MessageRepositoryController struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewMessageRepositoryController(db *sql.DB, logger *logrus.Logger) usecase.MessageRepository {
	return &MessageRepositoryController{
		db:     db,
		logger: logger,
	}
}

func (mrc *MessageRepositoryController) AddChatMessage(msg *mailer.Message) (*mailer.MessageInfo, error) {
	crMsg := &mailer.MessageInfo{}
	err := mrc.db.QueryRow(`INSERT INTO msg (author_id, chat_id, content, created_at) VALUES ($1, $2, $3, $4) 
	RETURNING message_id, author_id, chat_id, content, created_at`,
		msg.SenderId, msg.ChatId, msg.Content, msg.CreatedAt).Scan(&crMsg.Id, &crMsg.SenderId, &crMsg.ChatId, &crMsg.Content, &crMsg.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("psql CreateMessage: %w", err)
	}

	mrc.logger.WithField("message was succesfully created with messageID", crMsg.Id).Info("createMessage func")
	return crMsg, nil
}

func (mrc *MessageRepositoryController) DeleteMessage(messageID uint64) error {
	_, err := mrc.db.Exec(`DELETE FROM msg WHERE message_id=$1`, messageID)

	if err != nil {
		return fmt.Errorf("psql DeleteMessage: %w", err)
	}

	mrc.logger.WithField("message was successfully deleted with messageID", messageID).Info("deleteMessage func")
	return nil
}

func (mrc *MessageRepositoryController) UpdateMessage(msg *models.MessageUpdate) error {
	_, err := mrc.db.Exec(`UPDATE msg SET content=$1 WHERE message_id=$2`, msg.Content, msg.ID)

	if err != nil {
		return fmt.Errorf("psql UpdateMessage: %w", err)
	}

	mrc.logger.WithField("message was successfully updated with messageID", msg.ID).Info("updateMessage func")
	return nil

}

func (mrc *MessageRepositoryController) GetChatMessages(chatID uint64) ([]*mailer.MessageInfo, error) {
	rows, err := mrc.db.Query(`SELECT message_id, chat_id, author_id, content, created_at FROM msg WHERE chat_id=$1`, chatID)
	if err != nil {
		return nil, fmt.Errorf("getChatMessages: %w", err)
	}
	defer rows.Close()

	var messageList []*mailer.MessageInfo
	for rows.Next() {
		message := &mailer.MessageInfo{}
		if err := rows.Scan(&message.Id,
			&message.ChatId,
			&message.SenderId,
			&message.Content,
			&message.CreatedAt); err != nil {
			return nil, fmt.Errorf("getChatMessages rows.Next: %w", err)
		}
		messageList = append(messageList, message)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getChatMessages rows.Err: %w", err)
	}
	return messageList, nil
}

func (mrc *MessageRepositoryController) CreateChat() (*mailer.ChatInfo, error) {
	var chatID uint64
	err := mrc.db.QueryRow(`INSERT INTO chat DEFAULT VALUES RETURNING chat_id`).Scan(&chatID)

	if err != nil {
		return nil, fmt.Errorf("psql CreateChat: %w", err)
	}

	mrc.logger.WithField("chat was succesfully created with chatID", chatID).Info("createChat func")
	return &mailer.ChatInfo{ChatId: chatID}, nil
}

func (mrc *MessageRepositoryController) AddUserToChat(chatID uint64, userID uint64) error {
	var createdChatID, createdUserID uint64
	err := mrc.db.QueryRow(`INSERT INTO user_chat (user_id, chat_id) VALUES ($1, $2)
	 RETURNING user_id, chat_id`, userID, chatID).Scan(&createdUserID, &createdChatID)

	if err != nil {
		return fmt.Errorf("psql AddUserToChat: %w", err)
	}
	mrc.logger.WithField("user successfully added to chat", createdUserID).Info("addUserToChat func")
	return nil
}

func (mrc *MessageRepositoryController) GetChatUsers(chatID uint64) ([]*mailer.User, error) {
	rows, err := mrc.db.Query(`SELECT user_id FROM user_chat WHERE chat_id=$1`, chatID)

	if err != nil {
		return nil, fmt.Errorf("psql GetChatUsers %w", err)
	}
	defer rows.Close()

	users := make([]*mailer.User, 0)
	for rows.Next() {
		user := &mailer.User{}
		if err := rows.Scan(&user.Id); err != nil {
			return nil, fmt.Errorf("psql GetChatUsers rows.Next: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetChatUsers rows.Err: %w", err)
	}
	return users, nil
}

func (mrc *MessageRepositoryController) GetUserChats(userID uint64) ([]uint64, error) {
	rows, err := mrc.db.Query(`SELECT chat_id FROM user_chat WHERE user_id=$1`, userID)

	if err != nil {
		return nil, fmt.Errorf("psql GetUserChats %w", err)
	}
	defer rows.Close()

	chatIDs := make([]uint64, 0)
	for rows.Next() {
		var chatID uint64
		if err := rows.Scan(&chatID); err != nil {
			return nil, fmt.Errorf("psql GetUserChats rows.Next: %w", err)
		}
		chatIDs = append(chatIDs, chatID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUserChats rows.Err: %w", err)
	}
	fmt.Println("repository chat IDs", chatIDs)
	return chatIDs, nil
}

func (mrc *MessageRepositoryController) DeleteChat(chatID uint64) error {
	_, err := mrc.db.Exec(`DELETE FROM chat WHERE chat_id=$1`, chatID)

	if err != nil {
		return fmt.Errorf("psql DeleteChat: %w", err)
	}

	mrc.logger.WithField("chat with was successfully deleted with chatID", chatID).Info("deleteChat func")
	return nil
}
