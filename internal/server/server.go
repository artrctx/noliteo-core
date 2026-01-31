package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/artrctx/quoin-core/internal/config"
	"github.com/artrctx/quoin-core/internal/database"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	db   *database.Service
}

func NewServer() *http.Server {
	srvCfg := config.GetServerConfigFromEnv()
	srv := Server{srvCfg.Port, database.Get()}
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.port),
		Handler:      srv.Register(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
