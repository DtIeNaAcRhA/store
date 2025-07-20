package model

type User struct {
	ID           uint   `json:"id"`
	Username     string `json:"login"`
	HashPassword string `json:"password"`
}
