package dbclient

import (
	"github.com/deathcore666/battleShips/model"
)

type IGame interface {
	CreateGame(hostPlayer model.UserAccount) error
	JoinGame(guestPlayer model.UserAccount) error
}

type Game model.Game

func CreateGame(hostPlayer model.UserAccount) error {
	session, err := CreateSession(address, keyspace)
	if err != nil {
		return err
	}
	defer session.Close()

	iter := session.Query("SELECT * FROM games").Iter()
	var currentID = iter.NumRows() + 100000
	hostGame := Game{ID: currentID, P1: hostPlayer.UserName, P2: "", IsDone: false}

	query := "INSERT INTO games (id, p1, p2, isdone) VALUES (?, ?, ?, ?)"
	err = session.Query(query, hostGame.ID, hostGame.P1, hostGame.P2,
		hostGame.IsDone).Exec()
	return err
}

func JoinGame(guestPlayer model.UserAccount) error {
	session, err := CreateSession(address, keyspace)
	if err != nil {
		return err
	}
	defer session.Close()

	return err
}
