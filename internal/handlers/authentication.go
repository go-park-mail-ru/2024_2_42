package handlers

import (
	"fmt"
	"net/http"
)

func LogIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Succesful login")
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Succesful logout")
}

func IsAuthorized(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "User authentication")
}
