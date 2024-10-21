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

// User handlers layer
func InitializeUserDeliveryLayer(mux *http.ServeMux) {
	userRepo := repository.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo)
	userDelivery := NewUserDelivery(userUsecase)

	authRequiredRoutes := map[string]http.HandlerFunc{
		"POST /logout": userDelivery.LogOut,
	}

	authNotRequiredRoutes := map[string]http.HandlerFunc{
		"POST /login":        userDelivery.LogIn,
		"POST /signup":       userDelivery.SignUp,
		"GET /is_authorized": userDelivery.IsAuthorized,
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

// Feed handlers layer
func InitializeFeedDeliveryLayer(router *http.ServeMux) {
	feedRepo := repository.NewFeedRepository()
	feedUsecase := usecase.NewFeedUsecase(feedRepo)
	feedDelivery := NewFeedDelivery(feedUsecase)

	router.HandleFunc("/feed", feedDelivery.Feed)
}

// Routings handler
func Route() {
	routerParams := configs.NewInternalParams()
	mux := http.NewServeMux()

	InitializeUserDeliveryLayer(mux)
	InitializeFeedDeliveryLayer(mux)

	server := http.Server{
		Addr:    routerParams.MainServerPort,
		Handler: middleware.CORS(middleware.RequestID(middleware.Panic(mux))),
	}

	fmt.Printf("starting server at %s\n", routerParams.MainServerPort)
	log.Fatal(server.ListenAndServe())
}
