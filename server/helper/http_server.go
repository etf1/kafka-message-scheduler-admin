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

func ShutdownHTTPServer(srv *http.Server) error {
	defer log.Printf("http server shutted down")
	log.Printf("shutting down http server: %v", srv.Addr)

	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func StartupHTTPServer(srv *http.Server) {
	go func() {
		log.Printf("starting http server on %s", srv.Addr)
		defer log.Printf("http server stopped")

		err := srv.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				log.Printf("received closed event: %v", err)
			} else {
				log.Errorf("unexpected http server exit: %v", err)
			}
		}
	}()
}
