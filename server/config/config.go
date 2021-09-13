package config

import (
	"os"
	"strings"
	"time"

	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	log "github.com/sirupsen/logrus"
)

var (
	serverAddr             = ""
	DefaultShutdownTimeout = 5 * time.Second
)

func getBool(name string, defaultValue bool) bool {
	value, set := os.LookupEnv(name)
	if set {
		return value == "yes"
	}
	return defaultValue
}

func getString(name, defaultValue string) string {
	value, set := os.LookupEnv(name)
	if set {
		return value
	}
	return defaultValue
}

func getStrings(name string, defaultValue []string) []string {
	value, set := os.LookupEnv(name)
	if set {
		return helper.SplitTrim(value)
	}
	return defaultValue
}

func LogLevel() log.Level {
	lvl, err := log.ParseLevel(getString("LOG_LEVEL", "info"))
	if err != nil {
		return log.InfoLevel
	}
	return lvl
}

func GraylogServer() string {
	return getString("GRAYLOG_SERVER", "")
}

func MetricsAddr() string {
	return getString("METRICS_ADDR", ":9001")
}

// used for tests only
func SetServerAddr(addr string) string {
	serverAddr = addr
	return addr
}

func ServerAddr() string {
	if serverAddr != "" {
		return serverAddr
	}
	return getString("SERVER_ADDR", ":9000")
}

func SchedulersAddr() []string {
	return getStrings("SCHEDULERS_ADDR", []string{"localhost:8000"})
}

func StaticFilesDir() string {
	dir := getString("STATIC_FILES_DIR", "../client/build")
	return dir
}

func APIServerOnly() bool {
	return getBool("API_SERVER_ONLY", false)
}

func KafkaMessageBodyDecoder() string {
	return getString("KAFKA_MESSAGE_BODY_DECODER", "")
}

func DataRootDir() string {
	dir := getString("DATA_ROOT_DIR", "./.db")
	if !strings.HasSuffix(dir, "/") {
		return dir + "/"
	}
	return dir
}
