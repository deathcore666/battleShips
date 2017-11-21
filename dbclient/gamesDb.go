package dbclient

import (
	"github.com/deathcore666/battleShips/model"
)

type IGame interface {
	CreateGame(hostPlayer model.UserAccount) error
	JoinGame(guestPlayer model.UserAccount) error
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
