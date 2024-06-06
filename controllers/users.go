package controller

import (
	"net/http"

	"github.com/hanzalahimran7/go-auth/store"
)

type UserController struct {
	db store.DatabaseStore
}

func GetController(db store.DatabaseStore) *UserController {
	return &UserController{db}
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {}

func (uc *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {}

func (uc *UserController) DeleteProfile(w http.ResponseWriter, r *http.Request) {}

func (uc *UserController) EditProfile(w http.ResponseWriter, r *http.Request) {}
