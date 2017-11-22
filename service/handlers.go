package service

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/deathcore666/battleShips/dbclient"
	"github.com/deathcore666/battleShips/model"
	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func GetGamesHandler(w http.ResponseWriter, r *http.Request) {
	response, err := dbclient.GetGamesJSON()
	if err != nil {
		b, err := json.Marshal(err.Error())
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, string(b))
	}
	fmt.Fprint(w, string(response))
}

func JoinGameHandler(w http.ResponseWriter, r *http.Request) {
	gameID := r.FormValue("gameid")
	gameIDint, _ := strconv.Atoi(gameID)

	username := GetCookieField("name", "session", r)
	userID, err := dbclient.GetUserID(username)
	log.Println(gameID, username, userID)
	if err != nil {
		b, err := json.Marshal(err.Error())
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, string(b))
		return
	}

	err = dbclient.JoinGame(userID, gameIDint)
	if err != nil {
		b, err := json.Marshal(err.Error())
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, string(b))
		return
	}

	b, err := json.Marshal("joined successfuly")
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, string(b))
}

func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")

	username := GetCookieField("name", "session", r)
	userID, err := dbclient.GetUserID(username)
	log.Println(title, username, userID)
	if err != nil {
		b, err := json.Marshal(err.Error())
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, string(b))
		return
	}
	err = dbclient.CreateGame(userID, title)
	if err != nil {
		b, err := json.Marshal(err.Error())
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, string(b))
		return
	}
	b, err := json.Marshal("game created")
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, string(b))
	return
}

func IndexpageHandler(w http.ResponseWriter, r *http.Request) {
	if GetCookieField("name", "session", r) != "" {
		http.Redirect(w, r, "/start", 302)
	}

	errStr := GetCookieField("loginError", "errorCookie", r)
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
	var user model.UserAccount
	user.UserName = r.FormValue("name")
	user.Password = r.FormValue("password")
	redirectTarget := "/"
	if user.UserName != "" && user.Password != "" {
		err := dbclient.QueryUser(user)
		if err != nil {
			b, err := json.Marshal(err.Error())
			if err != nil {
				log.Println(err)
			}
			fmt.Fprintf(w, string(b))
			return
		}
		setSession("name", user.UserName, "session", w)
		resp := "000-success"
		b, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, string(b))
	}
	http.Redirect(w, r, redirectTarget, 302)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	userName := GetCookieField("name", "session", r)
	setSession("gameid", "123123", "gamesession", w)
	gameID := GetCookieField("gameid", "gamesession", r)
	log.Println(userName)
	b, err := json.Marshal(userName)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, string(b))
	c, err := json.Marshal(gameID)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, string(c))
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
	var user model.UserAccount
	user.UserName = r.FormValue("name")
	user.Password = r.FormValue("password")
	log.Println(r.Body)
	//redirectTarget := "/"
	if user.UserName != "" && user.Password != "" {
		err := dbclient.InsertUser(user)
		if err != nil {
			b, err := json.Marshal(err.Error())
			if err != nil {
				log.Println(err)
			}
			fmt.Fprintf(w, string(b))
			return
		}
	}
	if user.UserName == "" && user.Password == "" {
		fmt.Fprintln(w, "empty fields")
		fmt.Fprintln(w, r.Body)
		return
	}
	resp := "100-success"
	b, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, string(b))
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
