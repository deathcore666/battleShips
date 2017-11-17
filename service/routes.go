package service

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Indexpage",
		"GET",
		"/",
		IndexpageHandler,
	},
	Route{
		"Login",
		"POST",
		"/login",
		LoginHandler,
	},
	Route{
		"Logout",
		"POST",
		"/logout",
		LogoutHandler,
	},
	Route{
		"Startpage",
		"GET",
		"/start",
		StartpageHandler,
	},
}
