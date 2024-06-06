package goauth

import (
	"log"
	"net/http"
	"os"
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
	db := store.NewPostgresDB(
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("DB"),
	)
	if err := db.RunMigration(); err != nil {
		log.Fatal(err)
	}
	controller := controller.GetController(db)
	return &App{
		Router:     router,
		DB:         db,
		Controller: controller,
	}
}

func (a *App) Run() {
	a.LoadRoutes()
	http.ListenAndServe(":3000", a.Router)
}

func timeoutRequestMiddleware(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, 2*time.Second, "Request Timed Out")
}
