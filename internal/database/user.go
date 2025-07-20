package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"store/internal/model"
)

func CreateUser(user *model.User) error {
	stmt := `INSERT INTO user (username, hash_password) VALUES (?, ?)`
	res, err := DB.Exec(stmt, user.Username, user.HashPassword)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId() //что делает метод LastInsertId?
	user.ID = uint(id)
	return nil
}

func GetUserByLogin(login string) (*model.User, error) {
	row := DB.QueryRow(`SELECT id, username, hash_password FROM user WHERE username = ?`, login)
	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.HashPassword)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func GetUserByID(userID int) (*model.User, error) {
	log.Println(userID)
	row := DB.QueryRow(`SELECT id, username FROM user WHERE id = ?`, userID)
	log.Println(row)
	var user model.User
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not found")
			return nil, fmt.Errorf("user not found")

		}
		log.Println(err, "user not found")
		return nil, err

	}
	log.Println(user)
	return &user, nil
}
