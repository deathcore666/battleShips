package main

import (
	"net/http"

	"github.com/deathcore666/battleShips/service"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

func main() {
	router.HandleFunc("/", controller.IndexpageHandler)
	router.HandleFunc("/start", controller.StartpageHandler)

	router.HandleFunc("/login", controller.LoginHandler).Methods("POST")
	router.HandleFunc("/logout", controller.LogoutHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
