package dbclient

import (
	"errors"

	"github.com/gocql/gocql"
)

var (
	address  string = "127.0.0.1"
	keyspace string = "bships"
)

type ICassClient interface {
	CreateSession(address, keyspace string) (*gocql.Session, error)
	QueryUser(userName, password string) error
	InsertUser(userName, password string) error
}

func CreateSession(address, keyspace string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(address)
	cluster.Keyspace = keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func InsertUser(userName, password string) error {
	session, err := CreateSession(address, keyspace)
	if err != nil {
		return err
	}
	defer session.Close()
	query := "INSERT INTO users (lastname, age, city, email, firstname) VALUES ('Jones', 35, 'Austin', 'bob@example.com', 'Bob')"
	if err := session.Query(query).Exec(); err != nil {
		return err
	}
}

func QueryUser(userName, password string) error {
	session, err := CreateSession(address, keyspace)
	var pwd string
	if err != nil {
		return err
	}
	defer session.Close()

	iter := session.Query("SELECT password FROM users WHERE userName = (?)", userName).Iter()
	if iter.NumRows() == 0 {
		err := errors.New("user does not exist")
		return err
	}

	for iter.Scan(&pwd) {
		if pwd == password {
			return nil
		}
	}
	err = errors.New("wrong password")
	return err
}
