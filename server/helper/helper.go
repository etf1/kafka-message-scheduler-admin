package helper

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	lorem "github.com/drhodes/golorem"
)

func ShutdownHttpServer(srv *http.Server, defaultShutdownTimeout time.Duration) error {
	defer log.Printf("http server shut down")
	log.Printf("shutting down http server")

	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func Lipsum() string {
	return lorem.Paragraph(50, 100)
}

func SplitTrim(s string) []string {
	arr := strings.Split(s, ",")
	if len(arr) == 0 {
		return nil
	}
	result := make([]string, len(arr))
	for i := 0; i < len(arr); i++ {
		result[i] = strings.TrimSpace(arr[i])
	}
	return result
}

func VerifyIfSkipIntegrationTests(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "yes" {
		t.Skipf("skipping integration tests")
	}
}

// tells if the tests is running in docker
func IsRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); os.IsNotExist(err) {
		return false
	}
	return true
}

func GenRandString(prefix string) string {
	return fmt.Sprintf("%s%v", prefix, GenRandNum())
}

func GenRandNum() int64 {
	bg := big.NewInt(1000000)

	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}

	return n.Int64()
}

func Get(host, url string, timeout time.Duration) (*http.Response, error) {
	full := "http://" + host + url
	log.Debugf("calling get url: %v", full)

	req, err := http.NewRequest(http.MethodGet, full, nil)

	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: timeout,
	}

	return client.Do(req)
}

type CheckResponse func(*http.Response) error

var CheckResponseStatusOK = CheckResponse(func(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %v", resp.StatusCode)
	}
	return nil
})

var CheckResponseNil = CheckResponse(func(resp *http.Response) error {
	return nil
})

func DecodeJson(host, url string, timeout time.Duration, v interface{}, checkResponse ...CheckResponse) error {
	resp, err := Get(host, url, 5*time.Second)
	if err != nil {
		return fmt.Errorf("cannot get info from host %v: %v", host, err)
	}
	defer resp.Body.Close()

	check := CheckResponseStatusOK
	if len(checkResponse) > 0 {
		check = checkResponse[0]
	}

	err = check(resp)
	if err != nil {
		return fmt.Errorf("invalid response from host %v: %w", host, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot readall of the body from host %v: %v", host, resp.StatusCode)
	}
	log.Printf("http response body: %s", body)

	err = json.NewDecoder(bytes.NewReader(body)).Decode(v)
	if err != nil {
		return fmt.Errorf("unable to unmarshall json body: %v", err)
	}

	return nil
}
