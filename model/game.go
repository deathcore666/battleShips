package model

type Game struct {
	ID     string
	P1     UserAccount
	P2     UserAccount
	IsDone bool
}
