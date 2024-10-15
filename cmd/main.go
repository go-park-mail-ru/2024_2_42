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
	routerParams := configs.NewInternalParams()

	router := mux.NewRouter()

	// User Layer
	userRepo := repository.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo)
	userDelivery := delivery.NewUserDelivery(userUsecase, router)

	// Feed layer
	feedRepo := repository.NewFeedRepository()
	feedUsecase := usecase.NewFeedUsecase(feedRepo)
	feedDelivery := delivery.NewFeedDelivery(feedUsecase, router)

	// router.HandleFunc("/", handlers.Feed)
	// router.HandleFunc("/login", handlers.LogIn)
	// router.HandleFunc("/logout", handlers.LogOut)
	// router.HandleFunc("/is_authorized", handlers.IsAuthorized)
	// router.HandleFunc("/signup", handlers.SignUp)

	fmt.Printf("starting server at %s\n", routerParams.MainServerPort)
	log.Fatal(http.ListenAndServe(routerParams.MainServerPort, router))
}
