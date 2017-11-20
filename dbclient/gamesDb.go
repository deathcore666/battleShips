package dbclient

import "github.com/deathcore666/battleShips/model"

func CreateGame(game model.Game) error {
	session, err := CreateSession(address, keyspace)
	if err != nil {
		return err
	}
	defer session.Close()

	query := "INSERT INTO games (id, p1, p2, isdone) VALUES (?, ?, ?, ?)"
	err = session.Query(query, game.ID, game.P1, game.P2, false).Exec()
	return err
}
