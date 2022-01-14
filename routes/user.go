package routes

import (
	"api.ethscrow/controllers/user"
	"api.ethscrow/utils/session"
	"github.com/go-chi/chi/v5"
)

func UserRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/challenge", user.RequestChallenge)
	router.Post("/submit-challenge", user.SubmitChallenge)
	router.With(session.ProtectedRoute).Post("/logout", user.Logout)

	return router
}
