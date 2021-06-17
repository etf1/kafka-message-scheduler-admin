package restapi

import (
	"net/http"
	"time"

	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	"github.com/etf1/kafka-message-scheduler-admin/server/db"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
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
			Handler:      cors.AllowAll().Handler(initRouter(coldDB, liveDB, resv)),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
		coldBD: coldDB,
		liveDB: liveDB,
	}

	return srv
}

func NewRouter(coldDB db.DB, liveDB db.DB, resv schedulers.Resolver) http.Handler {
	return cors.AllowAll().Handler(initRouter(coldDB, liveDB, resv))
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
	return helper.ShutdownHttpServer(s.Server, config.ShutdownTimeout())
}
