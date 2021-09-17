package helper

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/sergi/go-diff/diffmatchpatch"
	log "github.com/sirupsen/logrus"
)

type MockedServer struct {
	s          *httptest.Server
	successful int
	failed     []string
}

// mockServerForQuery returns a mock server that only responds to a particular query string.
func MockServerForQuery(query string, code int, body string) *MockedServer {
	server := &MockedServer{}

	server.s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if query != "" && r.URL.RawQuery != query {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(query, r.URL.RawQuery, false)
			log.Printf("Query != Expected Query: %s", dmp.DiffPrettyText(diffs))
			server.failed = append(server.failed, r.URL.RawQuery)
			http.Error(w, "fail", 999)
			return
		}
		server.successful++

		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprint(w, body)
	}))

	return server
}

// Create a mock HTTP Server that will return a response with HTTP code and body.
func MockServer(code int, body string) *httptest.Server {
	serv := MockServerForQuery("", code, body)
	return serv.s
}
