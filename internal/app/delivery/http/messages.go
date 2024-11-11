package delivery

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"pinset/internal/app/models"
	internal_errors "pinset/internal/errors"
)

func (mdc *MessageDelieveryController) HandShake(w http.ResponseWriter, r *http.Request) {
	conn, err := mdc.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrBadRequest,
		})
		return
	}

	var chatJoiner models.ChatJoiner
	err = conn.ReadJSON(&chatJoiner)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrBadRequest,
		})
		return
	}

	newChatUser := &models.ChatUser{ID: chatJoiner.ID, ChatID: chatJoiner.ChatID, Connection: conn}

	mdc.Usecase.AddOnlineUser(newChatUser)

	go mdc.HandleConn(newChatUser)
}

func (mdc *MessageDelieveryController) GetAllChatMessages(w http.ResponseWriter, r *http.Request) {

	var user models.ChatJoiner

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrBadRequest,
		})
		return
	}

	messages, err := mdc.Usecase.GetChatMessages(user.ChatID)
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

		if err != nil {
			mdc.Logger.Printf("failed to read message %v", err)
			if errors.Is(err, net.ErrClosed) {
				return
			}
			user.Connection.WriteJSON(models.WebSocketResponse{Type: "error", Data: "websocket message bad request"})
			continue
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

		for _, reseiverID := range chatUserIDs {
			if mdc.Usecase.IsOnlineUser(reseiverID) {
				reseiver := mdc.Usecase.GetOnlineUser(reseiverID)
				reseiver.Connection.WriteJSON(models.WebSocketResponse{Type: "message", Data: messageInfo})

			}
		}

	}

}
