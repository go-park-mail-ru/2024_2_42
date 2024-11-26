package mediarepository

import (
	"context"
	"fmt"
	"pinset/internal/app/models"
	"pinset/internal/app/models/response"
	"pinset/mailer-service/mailer"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (mrc *MediaRepositoryController) CreateMessage(msg *models.Message) (*models.MessageCreateInfo, error) {

	req := &mailer.Message{
		SenderId:  msg.SenderID,
		ChatId:    msg.ChatID,
		Content:   msg.Content,
		CreatedAt: timestamppb.New(msg.CreatedAt),
	}

	res, err := mrc.mailerManager.AddChatMessage(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("gRPC CreateMessage: %w", err)
	}

	crMsg := &models.MessageCreateInfo{
		ID:        res.Id,
		SenderID:  res.SenderId,
		ChatID:    res.ChatId,
		Content:   res.Content,
		CreatedAt: res.CreatedAt.AsTime(),
	}
	mrc.logger.WithField("message was succesfully created with messageID", crMsg.ID).Info("createMessage func")
	return crMsg, nil
}

func (mrc *MediaRepositoryController) GetChatMessages(chatID uint64) ([]*models.MessageInfo, error) {
	req := &mailer.ChatRequest{
		ChatId: chatID,
	}
	res, err := mrc.mailerManager.GetAllChatMessages(context.Background(), req)

	if err != nil {
		return nil, fmt.Errorf("gRPC GetChatMessages: %w", err)
	}
	messageList := make([]*models.MessageInfo, 0)
	for _, message := range res.Messages {
		messageList = append(messageList, &models.MessageInfo{
			ID:        message.Id,
			SenderID:  message.SenderId,
			ChatID:    message.ChatId,
			Content:   message.Content,
			CreatedAt: message.CreatedAt.AsTime(),
		})
	}
	return messageList, nil
}

func (mrc *MediaRepositoryController) CreateChat(req *models.ChatCreateRequest) (*models.ChatInfo, error) {

	mailerRequest := &mailer.ChatCreateRequest{
		UserId:      req.UserID,
		CompanionId: req.CompanionID,
	}

	res, err := mrc.mailerManager.CreateChat(context.Background(), mailerRequest)
	if err != nil {
		return nil, fmt.Errorf("gRPC CreateChat: %w", err)
	}
	timeBirth := res.Companion.BirthTime.AsTime()
	companion := &response.UserProfileResponse{
		UserName:    res.Companion.UserName.Value,
		NickName:    res.Companion.NickName,
		Description: &res.Companion.Description.Value,
		BirthTime:   &timeBirth,
		Gender:      &res.Companion.Gender.Value,
		AvatarUrl:   &res.Companion.AvatarUrl.Value,
	}

	chatInfo := &models.ChatInfo{ChatID: res.ChatId, Companion: *companion}

	mrc.logger.WithField("chat was succesfully created with chatID", res.ChatId).Info("createChat func")
	return chatInfo, nil
}

func (mrc *MediaRepositoryController) GetChatUsers(chatID uint64) ([]uint64, error) {

	mailerRequest := &mailer.ChatRequest{
		ChatId: chatID,
	}

	res, err := mrc.mailerManager.GetChatUsers(context.Background(), mailerRequest)
	if err != nil {
		return nil, fmt.Errorf("gRPC GetChatUsers: %w", err)
	}

	userList := make([]uint64, 0)
	for _, user := range res.Users {
		userList = append(userList, user.Id)
	}

	return userList, nil
}

func (mrc *MediaRepositoryController) GetUserChats(userID uint64) ([]*models.ChatInfo, error) {

	mailerUser := &mailer.User{
		Id: userID,
	}

	res, err := mrc.mailerManager.GetUserChats(context.Background(), mailerUser)
	if err != nil {
		return nil, fmt.Errorf("gRPC GetChatUsers: %w", err)
	}
	fmt.Println("userchats in Method", res)

	chatList := make([]*models.ChatInfo, 0)
	for _, chat := range res.Chats {
		companion := &response.UserProfileResponse{
			UserName: chat.Companion.UserName.Value,
			NickName: chat.Companion.NickName,
		}
		if chat.Companion.UserName != nil {
			companion.UserName = chat.Companion.UserName.Value
		}
		birthTime := chat.Companion.BirthTime.AsTime()
		if chat.Companion.BirthTime != nil {
			companion.BirthTime = &birthTime
		}
		if chat.Companion.Gender != nil {
			companion.Gender = &chat.Companion.Gender.Value
		}
		if chat.Companion.AvatarUrl != nil {
			companion.AvatarUrl = &chat.Companion.AvatarUrl.Value
		}
		chatList = append(chatList, &models.ChatInfo{
			ChatID:    chat.ChatId,
			Companion: *companion,
		})
	}
	fmt.Println("chatList", chatList)
	return chatList, nil
}
