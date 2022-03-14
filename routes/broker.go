package routes

import (
	"api.ethscrow/controllers/broker"
	"api.ethscrow/utils/session"
	"github.com/go-chi/chi/v5"
)

func BrokerRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.With(session.ProtectedRoute).Get("/{PoolId}", broker.ConnectToPool)
	router.With(session.ProtectedRoute).Post("/create", broker.CreatePool)
	router.With(session.ProtectedRoute).Delete("/{PoolId}", broker.DeletePool)
	router.With(session.ProtectedRoute).Post("/{PoolId}", broker.UpdatePoolState)
	router.With(session.ProtectedRoute).Post("/{PoolId}/resolve", broker.ResolveConflict)
	router.With(session.ProtectedRoute).Get("/{PoolId}/accept", broker.AcceptPool)
	router.With(session.ProtectedRoute).Post("/{PoolId}/withdraw/generate", broker.GenerateTransaction)
	router.With(session.ProtectedRoute).Post("/{PoolId}/withdraw", broker.ProcessTransaction)

	return router
}
