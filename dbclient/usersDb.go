package dbclient

import (
	"errors"
	"log"

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

func GetUserID(username string) (int, error) {
	session, err := CreateSession(address, keyspace)
	if err != nil {
		return 0, err
	}
	defer session.Close()

	query := "SELECT id FROM users WHERE username = ? ALLOW FILTERING"
	iter := session.Query(query, username).Iter()

	var id int
	if !iter.Scan(&id) {
		return 0, errors.New("iteration error cassandra")
	}
	return id, nil
}

func InsertUser(user model.UserAccount) error {
	log.Println("attempting to create a user: ", user.UserName, user.Password)
	session, err := CreateSession(address, keyspace)
	if err != nil {
		return err
	}
	defer session.Close()

	checkQuery := "SELECT userName FROM users WHERE userName = ? ALLOW FILTERING"
	err = session.Query(checkQuery, user.UserName).Exec()
	if err != nil {
		return err
	}
	iter := session.Query(checkQuery, user.UserName).Iter()

	if iter.NumRows() == 0 {
		checkId := "SELECT * FROM users"
		iterid := session.Query(checkId).Iter()
		currentId := iterid.NumRows() + 10000
		query := "INSERT INTO users (id, userName, password) VALUES (?, ?, ?)"
		err = session.Query(query, currentId, user.UserName, user.Password).Exec()
		return err
	}
	return errors.New("101-username-already-exists")
}

func QueryUser(user model.UserAccount) error {
	session, err := CreateSession(address, keyspace)
	var pwd string
	if err != nil {
		return err
	}
	defer session.Close()

	iter := session.Query("SELECT password FROM users WHERE userName = ? ALLOW FILTERING", user.UserName).Iter()
	if iter.NumRows() == 0 {
		err := errors.New("001-wrong-username")
		return err
	}

	for iter.Scan(&pwd) {
		if pwd == user.Password {
			return nil
		}
	}
	err = errors.New("002-wrong-password")
	return err
}
