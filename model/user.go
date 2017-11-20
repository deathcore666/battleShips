package model

type UserAccount struct {
	userName string
	password string
}

type Game struct {
	ID string
	p1 UserAccount
	p2 UserAccount
}
