package dbclient

import (
	"github.com/deathcore666/battleShips/model"
)

type IGame interface {
	CreateGame(hostPlayer model.UserAccount) error
	JoinGame(guestPlayer model.UserAccount) error
}

type Game struct {
	ID     int
	P1     model.UserAccount
	P2     model.UserAccount
	IsDone bool
}

func (hostGame Game) CreateGame(hostPlayer model.UserAccount) error {
	session, err := CreateSession(address, keyspace)
	if err != nil {
		return err
	}
	defer session.Close()

	iter := session.Query("SELECT * FROM games").Iter()
	var currentID = iter.NumRows() + 10000
	hostGame = Game{ID: currentID, P1: hostPlayer, P2: model.UserAccount{}, IsDone: false}

	query := "INSERT INTO games (id, p1, p2, isdone) VALUES (?, ?, ?, ?)"
	err = session.Query(query, hostGame.ID, hostGame.P1, hostGame.P2,
		hostGame.IsDone).Exec()
	return err
}

func (hostGame Game) JoinGame(guestPlayer model.UserAccount) error {
	session, err := CreateSession(address, keyspace)
	if err != nil {
		return err
	}
	defer session.Close()

	query := "INSERT INTO games (id, p1, p2, isdone) VALUES (?, ?, ?, ?)"
	err = session.Query(query, hostGame.ID, hostGame.P1.UserName, guestPlayer.UserName, true).Exec()
	return err
}
