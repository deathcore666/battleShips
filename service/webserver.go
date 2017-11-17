package service

import "net/http"

func StartWebserver(port string) {
	r := NewRouter()
	http.Handle("/", r)
	http.ListenAndServe(":"+port, nil)
}
