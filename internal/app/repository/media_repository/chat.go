package mediarepository

import (
	"fmt"
	"pinset/internal/app/models"
)

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
	 RETURNING userID, chatID`, userID, chatID).Scan(&createdUserID, &createdChatID)

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
