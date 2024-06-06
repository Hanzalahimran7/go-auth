package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hanzalahimran7/go-auth/model"
	"github.com/hanzalahimran7/go-auth/store"
	"github.com/hanzalahimran7/go-auth/utils"
)

type UserController struct {
	db store.DatabaseStore
}

func GetController(db store.DatabaseStore) UserController {
	return UserController{db}
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) (int, error) { return 0, nil }

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) (int, error) { return 0, nil }

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) (int, error) {
	// The handler register a new user into the system
	userRequest := model.SignupRequest{}
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		return http.StatusBadRequest, fmt.Errorf("INVALID REQUEST BODY")
	}

	// Input validation
	if err := utils.ValidateUserRequest(userRequest); err != nil {
		return http.StatusBadRequest, err
	}

	// Check if email already exists
	if err := uc.db.FindByEmail(r.Context(), userRequest.Email); err == nil {
		return http.StatusBadRequest, fmt.Errorf("EMAIL ALREADY EXISTS")
	} else if err != nil && err != sql.ErrNoRows {
		return http.StatusBadRequest, fmt.Errorf("INTERNAL SERVER ERROR")
	}

	// Encrypt the user password
	encryptedPassword, err := utils.EncryptPassword(userRequest.Password)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("INTERNAL SERVER ERROR")
	}

	userId := uuid.New()
	// Create new user
	createdAt := time.Now().UTC()
	user := model.User{
		Id:        userId,
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Email:     userRequest.Email,
		Password:  encryptedPassword,
		CreatedAt: &createdAt,
	}

	// Save the user in DB
	if err := uc.db.Signup(r.Context(), &user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return http.StatusInternalServerError, fmt.Errorf("INTERNAL SERVER ERROR")
	}

	//Send the response to client
	utils.WriteJSON(w, http.StatusCreated, user)
	return 0, nil
}

func (uc *UserController) GetProfile(w http.ResponseWriter, r *http.Request) (int, error) {
	return 0, nil
}

func (uc *UserController) DeleteProfile(w http.ResponseWriter, r *http.Request) (int, error) {
	return 0, nil
}

func (uc *UserController) EditProfile(w http.ResponseWriter, r *http.Request) (int, error) {
	return 0, nil
}
