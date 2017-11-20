package service

import (
	"html/template"
	"net/http"

	"github.com/deathcore666/battleShips/dbclient"
	"github.com/gorilla/securecookie"
)

var Dbclient dbclient.ICassClient

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func IndexpageHandler(w http.ResponseWriter, r *http.Request) {
	if GetCookieField("name", r) != "" {
		http.Redirect(w, r, "/start", 302)
	}

	errStr := GetCookieField("error", r)
	errMap := map[string]string{
		"loginError": errStr,
	}
	t, _ := template.ParseFiles("views/login.html")
	t.Execute(w, errMap)
}

func StartpageHandler(w http.ResponseWriter, r *http.Request) {
	userName := GetCookieField("name", r)
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
	redirectTarget := "/"
	if name != "" && pass != "" {
		err := dbclient.QueryUser(name, pass)
		if err != nil {
			setSession("error", err.Error(), w)
			http.Redirect(w, r, redirectTarget, 302)
			return
		}
		clearSession(w)
		setSession("name", name, w)
		redirectTarget = "/start"
	}
	http.Redirect(w, r, redirectTarget, 302)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	redirectTarget := "/register"
	if name != "" && pass != "" {
		err := Dbclient.InsertUser(name, pass)
		if err != nil {
			redirectTarget = "/"
		}
	}
	http.Redirect(w, r, redirectTarget, 302)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

func setSession(field, fieldValue string, w http.ResponseWriter) {
	value := map[string]string{
		field: fieldValue,
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

func GetCookieField(field string, r *http.Request) (fieldValue string) {
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			fieldValue = cookieValue[field]
		}
	}
	return fieldValue
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
