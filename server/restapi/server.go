package restapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	"github.com/etf1/kafka-message-scheduler-admin/server/db"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers"
)

var (
	defaultShutdownTimeout = 5 * time.Second
)

type Server struct {
	*http.Server
	coldBD db.DB
	liveDB db.DB
}

func NewServer(coldDB db.DB, liveDB db.DB, resv schedulers.Resolver) Server {
	srv := Server{
		Server: &http.Server{
			Handler: cors.AllowAll().Handler(initRouter(coldDB, liveDB, resv)),
		},
		coldBD: coldDB,
		liveDB: liveDB,
	}

	return srv
}

func (s *Server) Router() http.Handler {
	return s.Handler
}

func (s *Server) Start(addr string) {
	go func() {
		log.Printf("starting rest api server on %s", addr)
		defer log.Printf("rest api server stopped")
		s.Addr = addr
		if err := s.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()
}

func (s *Server) Stop() error {
	defer log.Printf("rest api server shut down")
	log.Printf("shutting down rest api server")

	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
