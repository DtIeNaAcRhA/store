package handler

import (
	"encoding/json"
	"net/http"
	"regexp"
	"store/internal/config"
	"store/internal/database"
	"store/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)

	passwordRegex = regexp.MustCompile(`^.*([A-Z]|[А-ЯЁ]).*$`) // Хотя бы 1 заглавная
	digitRegex    = regexp.MustCompile(`^.*[0-9].*$`)          // Хотя бы 1 цифра
	specialRegex  = regexp.MustCompile(`^.*[!@#$%^&*].*$`)     // Хотя бы 1 спецсимвол
	spaceRegex    = regexp.MustCompile(`^.*\s.*$`)             // Запрет пробелов
)

func validatePassword(password string) string {
	if len(password) < 8 {
		return "Пароль должен быть не короче 8 символов"
	}
	if spaceRegex.MatchString(password) {
		return "Пароль не должен содержать пробелов и прочих невидимых символов"
	}
	if !passwordRegex.MatchString(password) {
		return "Пароль должен содержать хотя бы 1 заглавную букву"
	}
	if !digitRegex.MatchString(password) {
		return "Пароль должен содержать хотя бы 1 цифру"
	}
	if !specialRegex.MatchString(password) {
		return "Пароль должен содержать хотя бы 1 спецсимвол (!@#$%^&*)"
	}
	return ""
}

// POST/register
func Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSON(w, model.RegisterResponse{Message: "Некорректный запрос"}, http.StatusBadRequest)
		return
	}

	if !usernameRegex.MatchString(req.Username) {
		JSON(w, model.RegisterResponse{Message: "Логин должен содержать 3–20 символов и может содержать в себе: буквы, цифры, _"}, http.StatusBadRequest)
		return
	}

	if errMsg := validatePassword(req.Password); errMsg != "" {
		JSON(w, model.RegisterResponse{Message: errMsg}, http.StatusBadRequest)
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
