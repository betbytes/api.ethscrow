package routes

import (
	"api.ethscrow/controllers/broker"
	"github.com/go-chi/chi/v5"
)

func BrokerRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/connect/{roomId}", broker.Broker)

	return router
}
