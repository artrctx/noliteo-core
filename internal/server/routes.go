package server

import (
	"net/http"

	"github.com/artrctx/quoin-core/internal/server/middleware"
	"github.com/artrctx/quoin-core/internal/server/route"
	"github.com/artrctx/quoin-core/internal/service/token"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) Register() http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// health
	r.Get("/health", route.HealthHandlerFunc(s.db))

	// auth
	r.Get("/auth/verify", route.VerifyAuthHandler)

	// protected routes
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.Protected)

		// Token Route
		ts := token.TokenService{DB: s.db.Conn()}
		r.Get("/token/{token}", ts.ValidateTokenHandler)
	})

	return r
}
