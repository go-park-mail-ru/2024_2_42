package routing

import (
	"log"
	"net/http"
	"pinset/configs"

	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/middleware"
	"pinset/internal/app/repository"
	"pinset/internal/app/usecase"

	"pinset/pkg/logger"

	"github.com/sirupsen/logrus"
)

// Interfaces
type (
	UserDelivery interface {
		LogIn(w http.ResponseWriter, r *http.Request)
		LogOut(w http.ResponseWriter, r *http.Request)
		SignUp(w http.ResponseWriter, r *http.Request)
		IsAuthorized(w http.ResponseWriter, r *http.Request)
	}

	FeedDelivery interface {
		Feed(w http.ResponseWriter, r *http.Request)
	}

	MediaDelivery interface {
		GetMedia(w http.ResponseWriter, r *http.Request)
		UploadMedia(w http.ResponseWriter, r *http.Request)
	}
)

// Routings Main Handler
type RoutingHandler struct {
	logger      *logrus.Logger
	mux         *http.ServeMux
	userUsecase delivery.UserUsecase
}

func NewRoutingHandler(logger *logrus.Logger, mux *http.ServeMux, userUsecase delivery.UserUsecase) *RoutingHandler {
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
	authRequiredRoutes := map[string]http.HandlerFunc{
		"POST /logout": userHandlers.LogOut,
	}

	authNotRequiredRoutes := map[string]http.HandlerFunc{
		"POST /login":        userHandlers.LogIn,
		"POST /signup":       userHandlers.SignUp,
		"GET /is_authorized": userHandlers.IsAuthorized,
	}

	for route, handler := range authRequiredRoutes {
		rh.mux.HandleFunc(route, middleware.RequiredAuthorization(rh.logger, rh.userUsecase, handler))
	}

	for route, handler := range authNotRequiredRoutes {
		rh.mux.HandleFunc(route, middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, handler))
	}
}

func NewFeedDelivery(logger *logrus.Logger, usecase delivery.FeedUsecase) FeedDelivery {
	return &delivery.FeedDeliveryController{
		Usecase: usecase,
		Logger:  logger,
	}
}

// Feed layer handlers
func InitializeFeedLayerRoutings(rh *RoutingHandler, feedHandlers FeedDelivery) {
	rh.mux.HandleFunc("GET /feed", middleware.NotRequiredAuthorization(rh.logger, rh.userUsecase, feedHandlers.Feed))
}

func NewMediaDelivery(logger *logrus.Logger, usecase delivery.MediaUsecase) MediaDelivery {
	return &delivery.MediaDeliveryController{
		Usecase: usecase,
		Logger:  logger,
	}
}

// Media layer handlers
func InitializeMediaLayerRoutings(rh *RoutingHandler, mediaHandlers MediaDelivery) {
	rh.mux.HandleFunc("POST /create-pin", middleware.RequiredAuthorization(rh.logger, rh.userUsecase, mediaHandlers.UploadMedia))
}

func Route() {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	routerParams := configs.NewInternalParams()
	mux := http.NewServeMux()

	userRepo := repository.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo)
	userDelivery := NewUserDelivery(logger, userUsecase)

	feedRepo := repository.NewFeedRepository()
	feedUsecase := usecase.NewFeedUsecase(feedRepo)
	feedDelivery := NewFeedDelivery(logger, feedUsecase)

	mediaRepo, mediaErr := repository.NewMediaRepository()
	if mediaErr != nil {
		logger.Fatal(mediaErr)
	}
	mediaUsecase := usecase.NewMediaUsecase(mediaRepo)
	mediaDelivery := NewMediaDelivery(logger, mediaUsecase)

	rh := NewRoutingHandler(logger, mux, userUsecase)

	// Layers initialization
	InitializeUserLayerRoutings(rh, userDelivery)
	InitializeFeedLayerRoutings(rh, feedDelivery)
	InitializeMediaLayerRoutings(rh, mediaDelivery)

	server := http.Server{
		Addr:    routerParams.MainServerPort,
		Handler: middleware.CORS(middleware.RequestID(middleware.Panic(logger, mux))),
	}

	logger.WithField("starting server at ", routerParams.MainServerPort).Info()
	logger.Fatal(server.ListenAndServe())
}
