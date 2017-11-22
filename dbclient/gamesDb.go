package dbclient

import (
	"encoding/json"

	"github.com/deathcore666/battleShips/model"
)

type IGame interface {
	GetGamesJSON() ([]byte, error)
	CreateGame(hostPlayer model.UserAccount) error
	JoinGame(guestPlayer model.UserAccount) error
}

type GamesListJSON struct {
	Games []model.Game `json:"games"`
}

func GetGamesJSON() ([]byte, error) {
	session, err := CreateSession(address, keyspace)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var p1, id int
	var title string

	iter := session.Query("SELECT title, p1, id FROM games WHERE isdone = false ALLOW FILTERING").Iter()

	gamesList := make([]model.Game, iter.NumRows())
	for iter.Scan(&title, &p1, &id) {
		gamesList = append(gamesList, model.Game{ID: id, P1: p1, GameTitle: title})
	}

	res, err := json.Marshal(&GamesListJSON{Games: gamesList})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func CreateGame(hostPlayer model.UserAccount, ttle string) error {
	session, err := CreateSession(address, keyspace)
	if err != nil {
		return err
	}
	defer session.Close()

	iter := session.Query("SELECT * FROM games").Iter()
	var currentID = iter.NumRows() + 100000
	hostGame := model.Game{
		ID:        currentID,
		P1:        hostPlayer.ID,
		P2:        0,
		IsDone:    false,
		GameTitle: ttle,
	}

	query := "INSERT INTO games (id, p1, p2, isdone, title) VALUES (?, ?, ?, ?, ?)"
	err = session.Query(query, hostGame.ID, hostGame.P1, hostGame.P2,
		hostGame.IsDone, hostGame.GameTitle).Exec()
	return err
}

func JoinGame(guestPlayer model.UserAccount, gameID int) error {
	session, err := CreateSession(address, keyspace)
	if err != nil {
		return err
	}
	defer session.Close()

	query := "UPDATE games SET p2 = ?, isdone = ? WHERE id = ?"
	err = session.Query(query, guestPlayer.ID, true, gameID).Exec()
	return err
}
