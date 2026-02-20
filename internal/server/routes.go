package server

import (
	"net/http"

	"github.com/artrctx/noliteo-core/internal/hub"
	"github.com/artrctx/noliteo-core/internal/middleware"
	"github.com/artrctx/noliteo-core/internal/service/health"
	"github.com/artrctx/noliteo-core/internal/service/token"
	"github.com/artrctx/noliteo-core/internal/service/walktalk"
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
	r.Get("/health", health.HealthHandlerFunc(s.db))

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Token Route
	ts := token.TokenService{DB: s.db.Conn()}
	r.Post("/token", ts.GenerateTokenHandler)
	r.Get("/token", ts.ValidateTokenHandler)

	// protected routes
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.Protected)

		// Walkie Talkie Routes
		hub := hub.New()
		go hub.Run()
		wt := walktalk.WalkTalkService{DB: s.db.Conn(), Hub: hub}
		r.Get("/walkie-talkie", wt.WSHandler)
	})

	return r
}
