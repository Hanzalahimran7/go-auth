package goauth

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	controller "github.com/hanzalahimran7/go-auth/controllers"
	"github.com/hanzalahimran7/go-auth/store"
)

// App struct
type App struct {
	Router     *chi.Mux
	DB         store.DatabaseStore
	Controller controller.UserController
}

func Initialise() *App {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(timeoutRequestMiddleware)
	var db store.DatabaseStore = nil
	var c controller.UserController
	return &App{
		Router:     router,
		DB:         db,
		Controller: c,
	}
}

func (a *App) IntialiseDb(s store.DatabaseStore) {
	a.DB = s
	a.Controller = controller.GetController(s)
}

func (a *App) Run() {
	a.LoadRoutes()
	http.ListenAndServe(":3000", a.Router)
}

func timeoutRequestMiddleware(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, 2*time.Second, "Request Timed Out")
}
