package mediarepository

import (
	"fmt"
	"pinset/internal/app/models"
)

func (mrc *MediaRepositoryController) CreateMessage(msg *models.Message) (*models.MessageCreateInfo, error) {
	crMsg := &models.MessageCreateInfo{}
	err := mrc.db.QueryRow(`INSERT INTO msg (author_id, chat_id, content, created_at) VALUES ($1, $2, $3, $4) 
	RETURNING message_id, author_id, chat_id, content, created_at`,
		msg.SenderID, msg.ChatID, msg.Content, msg.CreatedAt).Scan(&crMsg.ID, &crMsg.SenderID, &crMsg.ChatID, &crMsg.Content, &crMsg.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("psql CreateMessage: %w", err)
	}

	mrc.logger.WithField("message was succesfully created with messageID", crMsg.ID).Info("createMessage func")
	return crMsg, nil
}

func (mrc *MediaRepositoryController) DeleteMessage(messageID uint64) error {
	_, err := mrc.db.Exec(`DELETE FROM msg WHERE message_id=$1`, messageID)

	if err != nil {
		return fmt.Errorf("psql DeleteMessage: %w", err)
	}

	mrc.logger.WithField("message was successfully deleted with messageID", messageID).Info("deleteMessage func")
	return nil
}

func (mrc *MediaRepositoryController) UpdateMessage(msg *models.MessageUpdate) error {
	_, err := mrc.db.Exec(`UPDATE msg SET content=$1 WHERE message_id=$2`, msg.Content, msg.ID)

	if err != nil {
		return fmt.Errorf("psql UpdateMessage: %w", err)
	}

	mrc.logger.WithField("message was successfully updated with messageID", msg.ID).Info("updateMessage func")
	return nil

}

func (mrc *MediaRepositoryController) GetChatMessages(chatID uint64) ([]*models.MessageInfo, error) {
	rows, err := mrc.db.Query(`SELECT message_id, chat_id, author_id, content, created_at FROM msg WHERE chat_id=$1`, chatID)
	if err != nil {
		return nil, fmt.Errorf("getChatMessages: %w", err)
	}
	defer rows.Close()

	var messageList []*models.MessageInfo
	for rows.Next() {
		message := &models.MessageInfo{}
		if err := rows.Scan(&message.ID,
			&message.ChatID,
			&message.SenderID,
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

func (mrc *MediaRepositoryController) CreateChat() (*models.ChatCreateInfo, error) {
	var chatID uint64
	err := mrc.db.QueryRow(`INSERT INTO chat DEFAULT VALUES RETURNING chat_id`).Scan(&chatID)

	if err != nil {
		return nil, fmt.Errorf("psql CreateChat: %w", err)
	}

	mrc.logger.WithField("chat was succesfully created with chatID", chatID).Info("createChat func")
	return &models.ChatCreateInfo{ID: chatID}, nil
}

func (mrc *MediaRepositoryController) AddUserToChat(chatID uint64, userID uint64) error {
	var createdChatID, createdUserID uint64
	err := mrc.db.QueryRow(`INSERT INTO user_chat (user_id, chat_id) VALUES ($1, $2)
	 RETURNING user_id, chat_id`, userID, chatID).Scan(&createdUserID, &createdChatID)

	if err != nil {
		return fmt.Errorf("psql AddUserToChat: %w", err)
	}
	mrc.logger.WithField("user successfully added to chat", createdUserID).Info("addUserToChat func")
	return nil
}

func (mrc *MediaRepositoryController) GetChatUsers(chatID uint64) ([]uint64, error) {
	rows, err := mrc.db.Query(`SELECT user_id FROM user_chat WHERE chat_id=$1`, chatID)

	if err != nil {
		return nil, fmt.Errorf("psql GetChatUsers %w", err)
	}
	defer rows.Close()

	userIDs := make([]uint64, 0)
	for rows.Next() {
		var userID uint64
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("psql GetChatUsers rows.Next: %w", err)
		}
		userIDs = append(userIDs, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetChatUsers rows.Err: %w", err)
	}
	return userIDs, nil
}

func (mrc *MediaRepositoryController) GetUserChats(userID uint64) ([]uint64, error) {
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

func (mrc *MediaRepositoryController) DeleteChat(chatID uint64) error {
	_, err := mrc.db.Exec(`DELETE FROM chat WHERE chat_id=$1`, chatID)

	if err != nil {
		return fmt.Errorf("psql DeleteChat: %w", err)
	}

	mrc.logger.WithField("chat with was successfully deleted with chatID", chatID).Info("deleteChat func")
	return nil
}
