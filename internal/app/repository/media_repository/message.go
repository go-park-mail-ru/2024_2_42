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
