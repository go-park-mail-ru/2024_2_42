package routing

import (
	"log"
	"net/http"
	"pinset/configs"

	"pinset/internal/app/db"
	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/middleware"
	mediarepository "pinset/internal/app/repository/media_repository"
	userRepository "pinset/internal/app/repository/user_repository"
	"pinset/internal/app/usecase"

	"pinset/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Interfaces
type (
	UserDelivery interface {
		LogIn(w http.ResponseWriter, r *http.Request)
		LogOut(w http.ResponseWriter, r *http.Request)
		SignUp(w http.ResponseWriter, r *http.Request)
		IsAuthorized(w http.ResponseWriter, r *http.Request)
		GetUserInfo(w http.ResponseWriter, r *http.Request)
		UpdateUserInfo(w http.ResponseWriter, r *http.Request)
	}

	MediaDelivery interface {
		Feed(w http.ResponseWriter, r *http.Request)

		GetPinPreview(w http.ResponseWriter, r *http.Request)
		GetPinPage(w http.ResponseWriter, r *http.Request)
		CreatePin(w http.ResponseWriter, r *http.Request)
		UpdatePin(w http.ResponseWriter, r *http.Request)
		DeletePin(w http.ResponseWriter, r *http.Request)

		GetUserBoards(w http.ResponseWriter, r *http.Request)
		GetBoard(w http.ResponseWriter, r *http.Request)
		CreateBoard(w http.ResponseWriter, r *http.Request)
		UpdateBoard(w http.ResponseWriter, r *http.Request)
		DeleteBoard(w http.ResponseWriter, r *http.Request)

		GetBookmark(w http.ResponseWriter, r *http.Request)
		CreateBookmark(w http.ResponseWriter, r *http.Request)
		DeleteBookmark(w http.ResponseWriter, r *http.Request)
	}
)

// Routings Main Handler
type RoutingHandler struct {
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
	rh.mux.HandleFunc("PUT /user/update/{user_id}", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.UpdateUserInfo))
	rh.mux.HandleFunc("POST /logout", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.LogOut))

	rh.mux.HandleFunc("POST /login", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.LogIn))
	rh.mux.HandleFunc("POST /signup", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.SignUp))
	rh.mux.HandleFunc("GET /is_authorized", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.IsAuthorized))
	rh.mux.HandleFunc("GET /user/{user_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.GetUserInfo))
}

func NewMediaDelivery(logger *logrus.Logger, usecase delivery.MediaUsecase) MediaDelivery {
	return &delivery.MediaDeliveryController{
		Usecase: usecase,
		Logger:  logger,
	}
}

// Media layer handlers
func InitializeMediaLayerRoutings(rh *RoutingHandler, mediaHandlers MediaDelivery) {
	rh.mux.HandleFunc("GET /feed", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.Feed))
	rh.mux.HandleFunc("POST /create-pin", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.CreatePin))

	rh.mux.HandleFunc("POST /create-pin", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.CreatePin))
	rh.mux.HandleFunc("GET /pins/preview/{pin_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.GetPinPreview))
	rh.mux.HandleFunc("GET /pins/page/{pin_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.GetPinPage))
	rh.mux.HandleFunc("PUT /pins/update/{pin_id}", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.UpdatePin))
	rh.mux.HandleFunc("DELETE /pins/delete/{pin_id}", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.DeletePin))

	rh.mux.HandleFunc("POST /create-board", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.CreateBoard))
	rh.mux.HandleFunc("GET /boards/{user_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.GetUserBoards))
	rh.mux.HandleFunc("GET /boards/{board_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.GetBoard))
	rh.mux.HandleFunc("PUT /boards/update/{board_id}", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.UpdateBoard))
	rh.mux.HandleFunc("DELETE /boards/delete/{board_id}", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.DeleteBoard))

	rh.mux.HandleFunc("POST /create-bookmark", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.CreateBookmark))
	rh.mux.HandleFunc("GET /bookmark/{bookmark_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.GetBookmark))
	rh.mux.HandleFunc("DELETE /bookmark/delete/{bookmark_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.DeleteBookmark))
}

func Route() {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	routerParams := configs.NewInternalParams()
	mux := mux.NewRouter()

	repo := db.InitDB(logger)

	userRepo := userRepository.NewUserRepository(repo, logger)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userDelivery := NewUserDelivery(logger, userUsecase)

	mediaRepo, mediaErr := mediarepository.NewMediaRepository(repo, logger)
	if mediaErr != nil {
		logger.Fatal(mediaErr)
	}

	logger.Info("MinioRepo created succesful!")

	mediaUsecase := usecase.NewMediaUsecase(mediaRepo)
	mediaDelivery := NewMediaDelivery(logger, mediaUsecase)

	rh := NewRoutingHandler(logger, mux, userUsecase)

	// Layers initialization
	InitializeUserLayerRoutings(rh, userDelivery)
	InitializeMediaLayerRoutings(rh, mediaDelivery)

	server := http.Server{
		Addr:    routerParams.MainServerPort,
		Handler: middleware.AccessLog(logger, middleware.CORS(middleware.RequestID(middleware.Panic(logger, mux)))),
	}

	logger.WithField("starting server at ", routerParams.MainServerPort).Info()
	logger.Fatal(server.ListenAndServe())
}
