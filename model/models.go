package model

type Game struct {
	ID        int
	P1        int
	P2        int
	IsDone    bool
	GameTitle string
}

type UserAccount struct {
	ID       int
	UserName string
	Password string
}
