package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	log "github.com/sirupsen/logrus"
)

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

func getInt(name string, defaultValue int) int {
	value, set := os.LookupEnv(name)
	if !set {
		return defaultValue
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
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

func MetricsHTTPAddr() string {
	return getString("METRICS_HTTP_ADDR", ":9001")
}

func APIServerAddr() string {
	return getString("API_SERVER_ADDR", ":9000")
}

func GroupID() string {
	return getString("GROUP_ID", "scheduler-admin-cg")
}

func SessionTimeout() int {
	return getInt("SESSION_TIMEOUT", 6000)
}

func SchedulersAddr() []string {
	return getStrings("SCHEDULERS_ADDR", []string{"localhost:8000"})
}

func DataRootDir() string {
	dir := getString("DATA_ROOT_DIR", "./.db")
	if !strings.HasSuffix(dir, "/") {
		return dir + "/"
	}
	return dir
}
