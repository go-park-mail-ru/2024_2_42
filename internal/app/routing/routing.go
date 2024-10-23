package routing

import (
	"fmt"
	"log"
	"net/http"
	"pinset/configs"
	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/middleware"
	"pinset/internal/app/repository"
	"pinset/internal/app/usecase"
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
)

func NewUserDelivery(usecase delivery.UserUsecase) UserDelivery {
	return &delivery.UserDeliveryController{
		Usecase: usecase,
	}
}

// User layer handlers
func InitializeUserLayerRoutings(mux *http.ServeMux, userUsecase delivery.UserUsecase, userHandlers UserDelivery) {
	authRequiredRoutes := map[string]http.HandlerFunc{
		"POST /logout": userHandlers.LogOut,
	}

	authNotRequiredRoutes := map[string]http.HandlerFunc{
		"POST /login":        userHandlers.LogIn,
		"POST /signup":       userHandlers.SignUp,
		"GET /is_authorized": userHandlers.IsAuthorized,
	}

	for route, handler := range authRequiredRoutes {
		mux.HandleFunc(route, middleware.RequiredAuthorization(userUsecase, handler))
	}

	for route, handler := range authNotRequiredRoutes {
		mux.HandleFunc(route, middleware.NotRequiredAuthorization(userUsecase, handler))
	}
}

func NewFeedDelivery(usecase delivery.FeedUsecase) FeedDelivery {
	return &delivery.FeedDeliveryController{
		Usecase: usecase,
	}
}

func InitializeFeedLayerRoutings(mux *http.ServeMux, userUsecase delivery.UserUsecase, feedHandlers FeedDelivery) {
	mux.HandleFunc("/feed", middleware.NotRequiredAuthorization(userUsecase, feedHandlers.Feed))
}

// Routings handler
func Route() {
	routerParams := configs.NewInternalParams()
	mux := http.NewServeMux()

	userRepo := repository.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo)
	userDelivery := NewUserDelivery(userUsecase)
	InitializeUserLayerRoutings(mux, userUsecase, userDelivery)

	feedRepo := repository.NewFeedRepository()
	feedUsecase := usecase.NewFeedUsecase(feedRepo)
	feedDelivery := NewFeedDelivery(feedUsecase)
	InitializeFeedLayerRoutings(mux, userUsecase, feedDelivery)

	server := http.Server{
		Addr:    routerParams.MainServerPort,
		Handler: middleware.CORS(middleware.RequestID(middleware.Panic(mux))),
	}

	fmt.Printf("starting server at %s\n", routerParams.MainServerPort)
	log.Fatal(server.ListenAndServe())
}
