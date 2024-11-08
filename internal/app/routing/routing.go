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
		AddPinToBoard(w http.ResponseWriter, r *http.Request)
		GetBoardPins(w http.ResponseWriter, r *http.Request)

		GetBookmark(w http.ResponseWriter, r *http.Request)
		CreateBookmark(w http.ResponseWriter, r *http.Request)
		DeleteBookmark(w http.ResponseWriter, r *http.Request)
		UploadMedia(w http.ResponseWriter, r *http.Request)
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
	rh.mux.HandleFunc("/user/{user_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, userHandlers.GetUserInfo)).Methods("GET")
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

	rh.mux.HandleFunc("/create-pin", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.CreatePin)).Methods("POST")
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
	rh.mux.HandleFunc("/boards/{board_id}/pins", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.GetBoardPins)).Methods("GET")

	rh.mux.HandleFunc("/create-bookmark", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.CreateBookmark)).Methods("POST")
	rh.mux.HandleFunc("/bookmark/{bookmark_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.GetBookmark)).Methods("GET")
	rh.mux.HandleFunc("/bookmark/delete/{bookmark_id}", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.DeleteBookmark)).Methods("DELETE")
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
