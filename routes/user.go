package routes

import (
	"api.ethscrow/controllers/user"
	"api.ethscrow/utils/session"
	"github.com/go-chi/chi/v5"
)

func UserRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/create", user.CreateUser)
	router.Get("/challenge/{Username}", user.RequestChallenge)
	router.Post("/challenge/{Username}", user.SubmitChallenge)
	router.With(session.ProtectedRoute).Get("/{Username}", user.PublicKey)
	router.With(session.ProtectedRoute).Get("/pool", user.AllPools)
	router.With(session.ProtectedRoute).Post("/logout", user.Logout)

	return router
}
