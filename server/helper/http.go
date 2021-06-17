package helper

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	defaultShutdownTimeout = 5 * time.Second
)

func ShutdownHttpServer(srv *http.Server) error {
	defer log.Printf("http server shutted down")
	log.Printf("shutting down http server")

	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func StartupHttpServer(srv *http.Server) chan bool {
	exitChan := make(chan bool, 1)

	go func() {
		log.Printf("starting http server on %s", srv.Addr)
		defer log.Printf("http server stopped")

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Errorf("unexpected http server exit: %v", err)
			exitChan <- true
		}
	}()

	return exitChan
}
