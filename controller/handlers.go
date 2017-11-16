package controller

import (
	"html/template"
	"net/http"

	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func IndexpageHandler(w http.ResponseWriter, r *http.Request) {
	if GetUserName(r) != "" {
		http.Redirect(w, r, "/start", 302)
	}
	t, _ := template.ParseFiles("views/login.html")
	t.Execute(w, nil)
}

func StartpageHandler(w http.ResponseWriter, r *http.Request) {
	userName := GetUserName(r)
	if userName != "" {
		userMap := map[string]string{
			"userName": userName,
		}
		t, _ := template.ParseFiles("views/start.html")
		t.Execute(w, userMap)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	redirectTarfet := "/"
	if name != "" && pass != "" {
		//TODO check credentials
		setSession(name, w)
		redirectTarfet = "/start"
	}
	http.Redirect(w, r, redirectTarfet, 302)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

func setSession(userName string, w http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func GetUserName(r *http.Request) (userName string) {
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
