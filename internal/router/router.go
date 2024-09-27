package router

import (
	"fmt"
	"log"
	"net/http"
	"youpin/configs"
	"youpin/internal/handlers"

	"github.com/gorilla/mux"
)

func Router() {
	routerParams := configs.NewInternalParams()

	router := mux.NewRouter()

	router.HandleFunc("/", handlers.Feed)
	router.HandleFunc("/login", handlers.LogIn)
	router.HandleFunc("/logout", handlers.LogOut)
	router.HandleFunc("/is_authorized", handlers.IsAuthorized)
	router.HandleFunc("/signup", handlers.SignUp)

	fmt.Printf("starting server at %s\n", routerParams.MainServerPort)
	log.Fatal(http.ListenAndServe(routerParams.MainServerPort, router))
}
