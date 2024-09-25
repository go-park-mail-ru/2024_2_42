package handlers

import (
	"fmt"
	"net/http"
)

func Feed(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Pins list")
}
