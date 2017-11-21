package dbclient

import (
	"errors"

	"github.com/deathcore666/battleShips/model"
	"github.com/gocql/gocql"
)

var (
	address  string = "127.0.0.1"
	keyspace string = "bships"
)

func CreateSession(address, keyspace string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(address)
	cluster.Keyspace = keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func InsertUser(user model.UserAccount) error {
	session, err := CreateSession(address, keyspace)
	if err != nil {
		return err
	}
	defer session.Close()

	checkQuery := "SELECT userName FROM users WHERE userName = ?"
	iter := session.Query(checkQuery, user.UserName).Iter()
	if iter.NumRows() == 0 {
		iterid := session.Query("SELECT * FROM users").Iter()
		var currentID = iterid.NumRows() + 10000
		query := "INSERT INTO users (id, userName, password) VALUES (?, ?, ?)"
		err = session.Query(query, currentID, user.UserName, user.Password).Exec()
		return err
	}
<<<<<<< HEAD
	return errors.New("101")
=======
	return errors.New("101-username-already-exists")
>>>>>>> games
}

func QueryUser(user model.UserAccount) error {
	session, err := CreateSession(address, keyspace)
	var pwd string
	if err != nil {
		return err
	}
	defer session.Close()

	iter := session.Query("SELECT password FROM users WHERE userName = ?", user.UserName).Iter()
	if iter.NumRows() == 0 {
<<<<<<< HEAD
		err := errors.New("001")
=======
		err := errors.New("001-wrong-username")
>>>>>>> games
		return err
	}

	for iter.Scan(&pwd) {
		if pwd == user.Password {
			return nil
		}
	}
<<<<<<< HEAD
	err = errors.New("002")
=======
	err = errors.New("002-wrong-password")
>>>>>>> games
	return err
}
