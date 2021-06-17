package runner

import (
	"net/http"
	"os"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	"github.com/etf1/kafka-message-scheduler-admin/server/db"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers"
	"github.com/etf1/kafka-message-scheduler-admin/server/restapi"
	"github.com/gorilla/mux"
)

type spaFileSystem struct {
	root http.FileSystem
}

// For single page app, we want to server index.html when file doesn't exist
func (fs *spaFileSystem) Open(name string) (http.File, error) {
	f, err := fs.root.Open(name)
	if os.IsNotExist(err) {
		return fs.root.Open("index.html")
	}
	return f, err
}

// TODO: accept a http.Server instance as parameter of the runner, if none then use a default server
func NewServer(coldDB, liveDB, historyDB db.DB, resolver schedulers.Resolver) *http.Server {
	var router http.Handler

	if config.APIServerOnly() {
		router = restapi.NewRouter(coldDB, liveDB, historyDB, resolver)
	} else {
		r := mux.NewRouter().StrictSlash(true)
		r.PathPrefix("/api").Handler(http.StripPrefix("/api", restapi.NewRouter(coldDB, liveDB, historyDB, resolver)))
		r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(&spaFileSystem{http.Dir(config.StaticFilesDir())})))
		router = r
	}

	return &http.Server{
		Handler:      router,
		Addr:         config.ServerAddr(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
