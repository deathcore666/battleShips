package model

type Game struct {
	ID     int
	P1     UserAccount
	P2     UserAccount
	IsDone bool
}

type UserAccount struct {
	UserName string
	Password string
}
