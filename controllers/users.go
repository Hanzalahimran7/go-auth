package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	return UserController{
		db: db,
	}
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) (int, error) { return 0, nil }

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) (int, error) { return 0, nil }

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) (int, error) {
	// The handler register a new user into the system
	userRequest := model.SignupRequest{}
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		return http.StatusBadRequest, fmt.Errorf("INVALID REQUEST BODY")
	}
	log.Println("******----------------------*********")
	// Input validation
	if err := utils.ValidateUserRequest(userRequest); err != nil {
		log.Printf("Invalid signup request: %v\nThe request body is %v\n", err, userRequest)
		return http.StatusBadRequest, err
	}
	log.Printf("User with email %s is signing up\n", userRequest.Email)

	// Check if email already exists
	log.Printf("Checking if the email %s exists in DB\n", userRequest.Email)
	err := uc.db.FindByEmail(r.Context(), userRequest.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Email %s already exists in Database\n", userRequest.Email)
			return http.StatusBadRequest, err
		}
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

	jwtToken, err := utils.CreateJWToken(user)
	if err != nil {
		log.Printf("Error creating JWT for user %s : %v+\n", user.Email, err)
		return http.StatusInternalServerError, fmt.Errorf("INTERNAL SERVER ERROR")
	}
	refreshToken, err := utils.CreateRefreshToken(user)
	if err != nil {
		log.Printf("Error creating Refresh for user %s: %v+\n", user.Email, err)
		return http.StatusInternalServerError, fmt.Errorf("INTERNAL SERVER ERROR")
	}
	// Save the user in DB
	log.Printf("Writing user %s to Database\n", userRequest.Email)
	if err := uc.db.Signup(r.Context(), &user); err != nil {
		log.Printf("Failed to write user %s to Database\n", userRequest.Email)
		w.WriteHeader(http.StatusInternalServerError)
		return http.StatusInternalServerError, fmt.Errorf("INTERNAL SERVER ERROR")
	}
	log.Printf("User %s created with email %s\n", user.Id, userRequest.Email)
	if err := uc.db.StoreToken(r.Context(), refreshToken, user.Id, time.Now().Add(time.Minute*2)); err != nil {
		log.Printf("Failed to write token for %s to Database: %v+\n", userRequest.Email, err)
		return http.StatusInternalServerError, fmt.Errorf("INTERNAL SERVER ERROR")
	}
	//Send the response to client
	// Add JWT and Refresh Token to the auth header
	w.Header().Add("Authorization", jwtToken)
	w.Header().Add("Refresh", refreshToken)
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
