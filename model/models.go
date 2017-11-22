package model

type Game struct {
	ID        int    `json:"ID"`
	P1        int    `json:"P1"`
	P2        int    `json:"P2"`
	IsDone    bool   `json:"isdone"`
	GameTitle string `json:"title"`
}

type UserAccount struct {
	ID       int    `json:"ID"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
