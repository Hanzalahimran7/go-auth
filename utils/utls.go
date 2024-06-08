package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hanzalahimran7/go-auth/model"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type apiFunc func(w http.ResponseWriter, r *http.Request) (int, error)

type APIError struct {
	Error string `json:"error"`
}

func ApiFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := f(w, r)
		if err != nil {
			WriteJSON(w, status, APIError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func ValidateUserRequest(userRequest model.SignupRequest) error {
	if userRequest.FirstName == "" || userRequest.LastName == "" || userRequest.Email == "" || userRequest.Password == "" {
		return fmt.Errorf("all fields are required")
	}
	if !isValidEmail(userRequest.Email) {
		return fmt.Errorf("invalid email address")
	}
	if len(userRequest.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	return nil
}

func ValidateLoginRequest(userRequest model.LoginRequest) error {
	if userRequest.Email == "" || userRequest.Password == "" {
		return fmt.Errorf("all fields are required")
	}
	if !isValidEmail(userRequest.Email) {
		return fmt.Errorf("invalid email address")
	}
	if len(userRequest.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	return nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+\/=?^_` + `"()` + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	return emailRegex.MatchString(email)
}

func CreateJWToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    user.Id,
			"email": user.Email,
			"iat":   time.Now().Unix(),
			"exp":   time.Now().Add(time.Minute * 1).Unix(),
			"iss":   "go-auth",
		})

	secret := os.Getenv("SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CreateRefreshToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  user.Id,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Minute * 2).Unix(),
			"iss": "go-auth",
		})
	secret := os.Getenv("SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
