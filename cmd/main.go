package main

import (
	"fmt"
	"log"
	"net/http"
	"pinset/configs"
	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/repository"
	"pinset/internal/app/usecase"

	"github.com/gorilla/mux"
)

func main() {
	// User Layer
	userRepo := repository.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo)
	userDelivery := delivery.NewUserDelivery(userUsecase)

	// Feed layer
	feedRepo := repository.NewFeedRepository()
	feedUsecase := usecase.NewFeedUsecase(feedRepo)
	feedDelivery := delivery.NewFeedDelivery(feedUsecase)

	// Routings handler
	routerParams := configs.NewInternalParams()
	router := mux.NewRouter()

	router.HandleFunc("/", feedDelivery.Feed)
	router.HandleFunc("/login", userDelivery.LogIn)
	router.HandleFunc("/logout", userDelivery.LogOut)
	router.HandleFunc("/signup", userDelivery.SignUp)
	router.HandleFunc("/is_authorized", userDelivery.IsAuthorized)

	fmt.Printf("starting server at %s\n", routerParams.MainServerPort)
	log.Fatal(http.ListenAndServe(routerParams.MainServerPort, router))
}
