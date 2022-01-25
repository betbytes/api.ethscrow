package main

import (
	"api.ethscrow/routes"
	"api.ethscrow/utils"
	"api.ethscrow/utils/database"
	"api.ethscrow/utils/session"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
)

func createRoute() *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.Throttle(20), // due to database connection limit
	)

	router.Route("/", func(r chi.Router) {
		r.Mount("/broker", routes.BrokerRoutes())
		r.Mount("/user", routes.UserRoutes())
	})

	return router
}

func main() {
	if err := utils.SetParams(); err != nil {
		log.Println(err)
		return
	}

	if err := database.ConnectToDatabase(); err != nil {
		log.Println(err)
		return
	}

	session.InitSessionStore()
	router := createRoute()

	walkF := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s, %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkF); err != nil {
		log.Fatalf("Logging Error: %s", err.Error())
	}

	log.Fatal(http.ListenAndServe(":"+utils.PORT, router))
	//log.Fatal(http.ListenAndServeTLS(":"+utils.PORT, "server.cert", "server.key", router))
}
