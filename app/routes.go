package goauth

import "github.com/hanzalahimran7/go-auth/utils"

func (a *App) LoadRoutes() {
	a.Router.Get("/login", utils.ApiFunc(a.Controller.Login))
	a.Router.Get("/logout", utils.ApiFunc(a.Controller.Logout))
	a.Router.Post("/register", utils.ApiFunc(a.Controller.Register))
	a.Router.Get("/profile", utils.ApiFunc(a.Controller.GetProfile))
	a.Router.Delete("/delete", utils.ApiFunc(a.Controller.DeleteProfile))
	a.Router.Patch("/edit", utils.ApiFunc(a.Controller.EditProfile))
}
