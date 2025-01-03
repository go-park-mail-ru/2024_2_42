package routing

import (
	"log"
	"net/http"
	"pinset/configs"

	"pinset/internal/app/db"
	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/middleware"
	mediarepository "pinset/internal/app/repository/media_repository"
	UserOnlineRepository "pinset/internal/app/repository/user_online_repository"
	userRepository "pinset/internal/app/repository/user_repository"
	"pinset/internal/app/usecase"

	"pinset/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	UserDelivery  UserDelivery
	MediaDelivery MediaDelivery
}

// Interfaces
type (
	UserDelivery interface {
		LogIn(w http.ResponseWriter, r *http.Request)
		LogOut(w http.ResponseWriter, r *http.Request)
		SignUp(w http.ResponseWriter, r *http.Request)
		IsAuthorized(w http.ResponseWriter, r *http.Request)
		GetAvatar(w http.ResponseWriter, r *http.Request)
		GetUserInfo(w http.ResponseWriter, r *http.Request)
		UpdateUserInfo(w http.ResponseWriter, r *http.Request)
		GetUsersByParams(w http.ResponseWriter, r *http.Request)
	}

	MediaDelivery interface {
		Feed(w http.ResponseWriter, r *http.Request)

		GetPinPreview(w http.ResponseWriter, r *http.Request)
		GetPinPage(w http.ResponseWriter, r *http.Request)
		CreatePin(w http.ResponseWriter, r *http.Request)
		UpdatePin(w http.ResponseWriter, r *http.Request)
		DeletePin(w http.ResponseWriter, r *http.Request)
		ViewPin(w http.ResponseWriter, r *http.Request)

		GetUserBoards(w http.ResponseWriter, r *http.Request)
		GetBoard(w http.ResponseWriter, r *http.Request)
		CreateBoard(w http.ResponseWriter, r *http.Request)
		UpdateBoard(w http.ResponseWriter, r *http.Request)
		DeleteBoard(w http.ResponseWriter, r *http.Request)
		AddPinToBoard(w http.ResponseWriter, r *http.Request)
		DeletePinFromBoard(w http.ResponseWriter, r *http.Request)
		GetBoardPins(w http.ResponseWriter, r *http.Request)

		GetBookmark(w http.ResponseWriter, r *http.Request)
		CreateBookmark(w http.ResponseWriter, r *http.Request)
		DeleteBookmark(w http.ResponseWriter, r *http.Request)
		UploadMedia(w http.ResponseWriter, r *http.Request)
	}

	MessageDelivery interface {
		HandShake(w http.ResponseWriter, r *http.Request)
		GetAllChatMessages(w http.ResponseWriter, r *http.Request)
		GetUserChats(w http.ResponseWriter, r *http.Request)
		CreateChat(w http.ResponseWriter, r *http.Request)
	}
)

// Routings Main Handler
type RoutingHandler struct {
	TestLogger  *logrus.Logger
	Mux         *mux.Router
	UserUsecase delivery.UserUsecase
	logger      *logrus.Logger
	mux         *mux.Router
	userUsecase delivery.UserUsecase
}

func NewRoutingHandler(logger *logrus.Logger, mux *mux.Router, userUsecase delivery.UserUsecase) *RoutingHandler {
	return &RoutingHandler{
		logger:      logger,
		mux:         mux,
		userUsecase: userUsecase,
	}
}

func NewUserDelivery(logger *logrus.Logger, usecase delivery.UserUsecase) UserDelivery {
	return &delivery.UserDeliveryController{
		Usecase: usecase,
		Logger:  logger,
	}
}

// User layer handlers
func InitializeUserLayerRoutings(rh *RoutingHandler, userHandlers UserDelivery) {
	rh.mux.HandleFunc("/user/update/{user_id}", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.UpdateUserInfo)).Methods("PUT")
	rh.mux.HandleFunc("/logout", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.LogOut)).Methods("POST")

	rh.mux.HandleFunc("/login", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.LogIn)).Methods("POST")
	rh.mux.HandleFunc("/signup", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.SignUp)).Methods("POST")
	rh.mux.HandleFunc("/is_authorized", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.IsAuthorized)).Methods("GET")
	rh.mux.HandleFunc("/get_avatar", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.GetAvatar)).Methods("GET")
	rh.mux.HandleFunc("/user/{user_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.GetUserInfo)).Methods("GET")
	rh.mux.HandleFunc("/users/by/params", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.GetUsersByParams)).Methods("POST")
}

func NewMediaDelivery(logger *logrus.Logger, usecase delivery.MediaUsecase) MediaDelivery {
	return &delivery.MediaDeliveryController{
		Usecase: usecase,
		Logger:  logger,
	}
}

// Media layer handlers
func InitializeMediaLayerRoutings(rh *RoutingHandler, mediaHandlers MediaDelivery) {
	rh.mux.HandleFunc("/image/upload", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.UploadMedia)).Methods("POST")

	rh.mux.HandleFunc("/feed", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.Feed)).Methods("GET")

	rh.mux.HandleFunc("/create-pin", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.CreatePin)).Methods("POST")
	rh.mux.HandleFunc("/pins/view/{pin_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.ViewPin)).Methods("POST")
	rh.mux.HandleFunc("/pins/preview/{pin_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.GetPinPreview)).Methods("GET")
	rh.mux.HandleFunc("/pins/page/{pin_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.GetPinPage)).Methods("GET")
	rh.mux.HandleFunc("/pins/update/{pin_id}", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.UpdatePin)).Methods("PUT")
	rh.mux.HandleFunc("/pins/delete/{pin_id}", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.DeletePin)).Methods("DELETE")

	rh.mux.HandleFunc("/create-board", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.CreateBoard)).Methods("POST")
	rh.mux.HandleFunc("/boards/{user_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.GetUserBoards)).Methods("GET")
	rh.mux.HandleFunc("/boards/{board_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.GetBoard)).Methods("GET")
	rh.mux.HandleFunc("/boards/update/{board_id}", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.UpdateBoard)).Methods("PUT")
	rh.mux.HandleFunc("/boards/delete/{board_id}", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.DeleteBoard)).Methods("DELETE")

	rh.mux.HandleFunc("/boards/{board_id}/addpin/{pin_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.AddPinToBoard)).Methods("POST")
	rh.mux.HandleFunc("/boards/{board_id}/deletepin/{pin_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.DeletePinFromBoard)).Methods("DELETE")
	rh.mux.HandleFunc("/boards/{board_id}/pins", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.GetBoardPins)).Methods("GET")

	rh.mux.HandleFunc("/create-bookmark", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.CreateBookmark)).Methods("POST")
	rh.mux.HandleFunc("/bookmark/{bookmark_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.GetBookmark)).Methods("GET")
	rh.mux.HandleFunc("/bookmark/delete/{bookmark_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.DeleteBookmark)).Methods("DELETE")

	// rh.mux.HandleFunc("/handshake", delivery.HandShake).Methods("GET")
}

func NewMessageDelivery(logger *logrus.Logger, usecase delivery.MessageUsecase) MessageDelivery {
	return &delivery.MessageDelieveryController{
		Usecase: usecase,
		Logger:  logger,
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func InitializeMessageLayerRoutings(rh *RoutingHandler, messageHandlers MessageDelivery) {
	rh.mux.HandleFunc("/handshake", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, messageHandlers.HandShake)).Methods("GET")
	rh.mux.HandleFunc("/chat/{chat_id}/messages", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, messageHandlers.GetAllChatMessages)).Methods("GET")
	rh.mux.HandleFunc("/mychats", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, messageHandlers.GetUserChats)).Methods("GET")
	rh.mux.HandleFunc("/create/chat/{user_id}", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, messageHandlers.CreateChat)).Methods("POST")
}

func Route() {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	routerParams := configs.NewInternalParams()
	mux := mux.NewRouter()

	repo := db.InitDB(logger)
	mediaRepo, mediaErr := mediarepository.NewMediaRepository(repo, logger)
	if mediaErr != nil {
		logger.Fatal(mediaErr)
	}

	userRepo := userRepository.NewUserRepository(repo, logger)
	userUsecase := usecase.NewUserUsecase(userRepo, mediaRepo)
	userDelivery := NewUserDelivery(logger, userUsecase)

	mediaUsecase := usecase.NewMediaUsecase(mediaRepo, userRepo)
	mediaDelivery := NewMediaDelivery(logger, mediaUsecase)

	userOnlineRepo := UserOnlineRepository.NewUserOnlineRepository()
	messageUsecase := usecase.NewMessageUsecase(userOnlineRepo, mediaRepo, userRepo)
	messageDelivery := NewMessageDelivery(logger, messageUsecase)

	rh := NewRoutingHandler(logger, mux, userUsecase)

	// Layers initialization
	InitializeUserLayerRoutings(rh, userDelivery)
	InitializeMediaLayerRoutings(rh, mediaDelivery)
	InitializeMessageLayerRoutings(rh, messageDelivery)

	server := http.Server{
		Addr:    routerParams.MainServerPort,
		Handler: middleware.AccessLog(logger, middleware.CORS(middleware.RequestID(middleware.Panic(logger, mux)))),
	}

	logger.WithField("starting server at ", routerParams.MainServerPort).Info()
	logger.Fatal(server.ListenAndServe())
}
