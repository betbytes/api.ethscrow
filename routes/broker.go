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
	router.With(session.ProtectedRoute).Post("/{PoolId}", broker.UpdatePoolState)
	router.With(session.ProtectedRoute).Post("/{PoolId}/resolve", broker.ResolveConflict)
	router.With(session.ProtectedRoute).Get("/{PoolId}/accept", broker.AcceptPool)

	return router
}
