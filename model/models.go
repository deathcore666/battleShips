package model

type Game struct {
	ID     int
	P1     string
	P2     string
	IsDone bool
}

type UserAccount struct {
	UserName string
	Password string
}
