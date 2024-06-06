package goauth

func (a *App) LoadRoutes() {
	a.Router.Get("/login", a.Controller.Login)
	a.Router.Get("/logout", a.Controller.Logout)
	a.Router.Post("/register", a.Controller.Register)
	a.Router.Get("/profile", a.Controller.GetProfile)
	a.Router.Delete("/delete", a.Controller.DeleteProfile)
	a.Router.Patch("/edit", a.Controller.EditProfile)
}
