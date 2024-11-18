package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pinset/configs"
	"pinset/internal/app/models"
	internal_errors "pinset/internal/errors"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (mdc *MessageDelieveryController) HandShake(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handshake started")
	conn, err := mdc.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrBadRequest,
		})
		return
	}

	userID, ok := r.Context().Value(configs.UserIdKey).(uint64)
	if !ok {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: internal_errors.ErrUserIsNotAuthorized, Internal: internal_errors.ErrUserIsNotAuthorized,
		})
	}

	fmt.Println("last ID connected ", userID)
	fmt.Println("num users online ", mdc.Usecase.NumUsersOnline())

	newChatUser := &models.ChatUser{ID: userID, Connection: conn}

	mdc.Usecase.AddOnlineUser(newChatUser)

	go mdc.HandleConn(newChatUser)
}

func (mdc *MessageDelieveryController) GetUserChats(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(configs.UserIdKey).(uint64)
	if !ok {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: internal_errors.ErrUserIsNotAuthorized, Internal: internal_errors.ErrUserIsNotAuthorized,
		})
	}
	chats, err := mdc.Usecase.GetUserChats(userID)
	fmt.Println(chats)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(chats); err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
		return
	}

}

func (mdc *MessageDelieveryController) GetAllChatMessages(w http.ResponseWriter, r *http.Request) {

	chatIDStr := mux.Vars(r)["chat_id"]
	chatID, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	messages, err := mdc.Usecase.GetChatMessages(chatID)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
		return
	}
}

func (mdc *MessageDelieveryController) HandleConn(user *models.ChatUser) {
	defer mdc.Usecase.DeleteOnlineUser(user.ID)
	defer user.Connection.Close()

	for {
		var mes models.Message
		err := user.Connection.ReadJSON(&mes)
		mes.SenderID = user.ID
		mes.CreatedAt = time.Now()
		fmt.Println(mes)

		if err != nil {
			mdc.Logger.Printf("failed to read message %v", err)
			return
			// if errors.Is(err, net.ErrClosed) {
			// 	mdc.Logger.Printf("going away %v", err)
			// 	return
			// }
			// user.Connection.WriteJSON(models.WebSocketResponse{Type: "error", Data: "websocket message bad request"})
			// continue
		}

		messageInfo, err := mdc.Usecase.AddChatMessage(&mes)
		if err != nil {
			mdc.Logger.Printf("failed to add message to chat %v", err)
			user.Connection.WriteJSON(models.WebSocketResponse{Type: "error", Data: "failed to add message to chat"})
			continue
		}

		chatID := mes.ChatID
		chatUserIDs, err := mdc.Usecase.GetChatUsers(chatID)
		if err != nil {
			mdc.Logger.Printf("failed get chat users %v", err)
			user.Connection.WriteJSON(models.WebSocketResponse{Type: "error", Data: "failed to get chat users"})
			continue
		}

		fmt.Println("messageInfo", messageInfo)
		for _, reseiverID := range chatUserIDs {
			if mdc.Usecase.IsOnlineUser(reseiverID) {
				reseiver := mdc.Usecase.GetOnlineUser(reseiverID)
				reseiver.Connection.WriteJSON(models.WebSocketResponse{Type: "message", Data: messageInfo})

			}
		}

	}

}
