package handler

import (
	"encoding/json"
	"net/http"
	"store/internal/config"
	"store/internal/database"
	"store/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// POST/register
func Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSON(w, model.RegisterResponse{Message: "Некорректный запрос"}, http.StatusBadRequest)
		return
	}

	// Проверка на существование
	existingUser, _ := database.GetUserByLogin(req.Username)
	if existingUser != nil && existingUser.ID != 0 {
		JSON(w, model.RegisterResponse{Message: "Пользователь уже существует"}, http.StatusConflict)
		return
	}

	// Хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		JSON(w, model.RegisterResponse{Message: "Ошибка сервера"}, http.StatusInternalServerError)
		return
	}

	user := &model.User{
		Username:     req.Username,
		HashPassword: string(hash),
	}
	if err := database.CreateUser(user); err != nil {
		JSON(w, model.RegisterResponse{Message: "Ошибка при сохранении"}, http.StatusInternalServerError)
		return
	}

	JSON(w, model.RegisterResponse{Message: "Пользователь успешно зарегистрирован"}, http.StatusCreated)
}

// POST/login
func Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSON(w, model.LoginResponse{Token: ""}, http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByLogin(req.Username)
	if err != nil || user.ID == 0 {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(req.Password)); err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	// Генерация токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	JSON(w, model.LoginResponse{Token: tokenString}, http.StatusOK)
}
