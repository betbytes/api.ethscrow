package routes

import (
	"api.ethscrow/controllers/broker"
	"api.ethscrow/utils/session"
	"github.com/go-chi/chi/v5"
)

func BrokerRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.With(session.ProtectedRoute).Get("/{PoolId}", broker.ConnectToPool)
	router.Post("/create", broker.CreatePool)
	router.Delete("/{PoolId}", broker.DeletePool)

	return router
}
