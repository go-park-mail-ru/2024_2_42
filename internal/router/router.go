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
	router.HandleFunc("/logout/{id:[0-9]+}", handlers.LogOut)
	router.HandleFunc("/is_authorized", handlers.IsAuthorized)

	fmt.Println("starting server at :8080")
	log.Fatal(http.ListenAndServe(routerParams.MainServerPort, router))
}
