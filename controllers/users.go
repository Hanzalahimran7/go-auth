package controller

import (
	"encoding/json"
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

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {
	// The handler register a new user into the system
	userRequest := model.SignupRequest{}
	// Gets the signup request data from request's body and save it in userRequest
	// In case of errors, send a bad request status
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Encrypt the user password and throw bad request in case of errors
	encryptedPassword, err := utils.EncryptPassword(userRequest.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Create a new UUID for user
	userId := uuid.New()
	// Get current time for created at
	createdAt := time.Now().UTC()
	user := model.User{
		Id:        userId,
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Email:     userRequest.Email,
		Password:  encryptedPassword,
		CreatedAt: &createdAt,
	}
	if err := uc.db.Signup(r.Context(), &user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (uc *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {}

func (uc *UserController) DeleteProfile(w http.ResponseWriter, r *http.Request) {}

func (uc *UserController) EditProfile(w http.ResponseWriter, r *http.Request) {}
