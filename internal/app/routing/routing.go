package routing

import (
	"fmt"
	"log"
	"net/http"
	"pinset/configs"
	
	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/repository"
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
	}

	FeedDelivery interface {
		Feed(w http.ResponseWriter, r *http.Request)
	}

	MediaDelivery interface {
		GetMedia(w http.ResponseWriter, r *http.Request)
		UploadMedia(w http.ResponseWriter, r *http.Request)
	}
)

func NewUserDelivery(usecase delivery.UserUsecase) UserDelivery {
	return &delivery.UserDeliveryController{
		Usecase: usecase,
	}
}

// User handlers layer
func InitializeUserDeliveryLayer(router *mux.Router) {
	userRepo := repository.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo)
	userDelivery := NewUserDelivery(userUsecase)

	router.HandleFunc("/login", userDelivery.LogIn)
	router.HandleFunc("/logout", userDelivery.LogOut)
	router.HandleFunc("/signup", userDelivery.SignUp)
	router.HandleFunc("/is_authorized", userDelivery.IsAuthorized)
}

func NewFeedDelivery(usecase delivery.FeedUsecase) FeedDelivery {
	return &delivery.FeedDeliveryController{
		Usecase: usecase,
	}
}

// Feed handlers layer
func InitializeFeedDeliveryLayer(router *mux.Router) {
	feedRepo := repository.NewFeedRepository()
	feedUsecase := usecase.NewFeedUsecase(feedRepo)
	feedDelivery := NewFeedDelivery(feedUsecase)

	router.HandleFunc("/", feedDelivery.Feed)
}

func Route() {
	// Routings handler
	routerParams := configs.NewInternalParams()
	router := mux.NewRouter()

	InitializeUserDeliveryLayer(router)
	InitializeFeedDeliveryLayer(router)

	err := logger.Logger()

	if err != nil {
		log.Fatal(fmt.Errorf("Route: %w", err))
	}

	logrus.WithField("starting server at ", routerParams.MainServerPort).Info()
	log.Fatal(http.ListenAndServe(routerParams.MainServerPort, router))
}
