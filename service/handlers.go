package service

import (
	"html/template"
	"net/http"

	"github.com/deathcore666/battleShips/dbclient"
	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func IndexpageHandler(w http.ResponseWriter, r *http.Request) {
	if GetCookieField("name", "session", r) != "" {
		http.Redirect(w, r, "/start", 302)
	}

	errStr := GetCookieField("error", "errorCookie", r)
	errMap := map[string]string{
		"loginError": errStr,
	}
	t, _ := template.ParseFiles("views/login.html")
	t.Execute(w, errMap)
}

func StartpageHandler(w http.ResponseWriter, r *http.Request) {
	userName := GetCookieField("name", "session", r)
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
			setSession("error", err.Error(), "errorCookie", w)
			http.Redirect(w, r, redirectTarget, 302)
			return
		}
		clearSession("errorCookie", w)
		setSession("name", name, "session", w)
		redirectTarget = "/start"
	}
	http.Redirect(w, r, redirectTarget, 302)
}

func RegisterpageHandler(w http.ResponseWriter, r *http.Request) {
	errStr := GetCookieField("registerError", "errorCookie", r)
	errMap := map[string]string{
		"registerError": errStr,
	}
	t, _ := template.ParseFiles("views/register.html")
	t.Execute(w, errMap)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	redirectTarget := "/"
	if name != "" && pass != "" {
		err := dbclient.InsertUser(name, pass)
		if err != nil {
			setSession("registerError", err.Error(), "errorCookie", w)
			redirectTarget = "/registerp"
			http.Redirect(w, r, redirectTarget, 302)
			return
		}
	}
	clearSession("errorCookie", w)
	http.Redirect(w, r, redirectTarget, 302)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession("session", w)
	http.Redirect(w, r, "/", 302)
}

func setSession(field, fieldValue, cookieType string, w http.ResponseWriter) {
	value := map[string]string{
		field: fieldValue,
	}
	if encoded, err := cookieHandler.Encode(cookieType, value); err == nil {
		cookie := &http.Cookie{
			Name:  cookieType,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func GetCookieField(field, cookieType string, r *http.Request) (fieldValue string) {
	if cookie, err := r.Cookie(cookieType); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode(cookieType, cookie.Value, &cookieValue); err == nil {
			fieldValue = cookieValue[field]
		}
	}
	return fieldValue
}

func clearSession(cookieType string, w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   cookieType,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
