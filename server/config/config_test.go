package config_test

import (
	"os"
	"testing"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
)

// TestKafkaMessageBodyDecoder covers the M2 fix: empty env var must return "",
// not "http://" (which would pass the != "" check in kafka.go and create a broken decoder).
func TestKafkaMessageBodyDecoder(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		// empty string → "" (M2 fix: must NOT return "http://")
		{"empty", "", ""},
		// already has http:// prefix → unchanged
		{"http prefix", "http://decoder.internal:8080", "http://decoder.internal:8080"},
		// already has https:// prefix → unchanged
		{"https prefix", "https://decoder.internal:8080", "https://decoder.internal:8080"},
		// bare host → http:// prepended
		{"bare host", "decoder.internal:8080", "http://decoder.internal:8080"},
		// uppercase HTTP:// → treated as already having a scheme (strings.ToLower check)
		{"uppercase HTTP", "HTTP://decoder.internal:8080", "HTTP://decoder.internal:8080"},
		// uppercase HTTPS:// → treated as already having a scheme
		{"uppercase HTTPS", "HTTPS://decoder.internal:8080", "HTTPS://decoder.internal:8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("KAFKA_MESSAGE_BODY_DECODER", tt.envValue)
			got := config.KafkaMessageBodyDecoder()
			if got != tt.want {
				t.Errorf("KafkaMessageBodyDecoder() with KAFKA_MESSAGE_BODY_DECODER=%q: got %q, want %q",
					tt.envValue, got, tt.want)
			}
		})
	}
}

// TestKafkaMessageBodyDecoder_unset ensures an unset env var returns "".
func TestKafkaMessageBodyDecoder_unset(t *testing.T) {
	orig, wasSet := os.LookupEnv("KAFKA_MESSAGE_BODY_DECODER")
	if err := os.Unsetenv("KAFKA_MESSAGE_BODY_DECODER"); err != nil {
		t.Fatalf("failed to unset env: %v", err)
	}
	t.Cleanup(func() {
		if wasSet {
			os.Setenv("KAFKA_MESSAGE_BODY_DECODER", orig)
		}
	})

	got := config.KafkaMessageBodyDecoder()
	if got != "" {
		t.Errorf("KafkaMessageBodyDecoder() with unset env: got %q, want %q", got, "")
	}
}

// TestGetBool_viaAPIServerOnly covers the M1 fix: getBool now uses strconv.ParseBool.
// Old behavior: only "yes" returned true. New behavior: "true", "1", "TRUE", etc.
func TestGetBool_viaAPIServerOnly(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     bool
	}{
		// standard Go boolean strings (all valid via strconv.ParseBool)
		{"true", "true", true},
		{"TRUE", "TRUE", true},
		{"True", "True", true},
		{"1", "1", true},
		{"t", "t", true},
		{"T", "T", true},
		{"false", "false", false},
		{"FALSE", "FALSE", false},
		{"False", "False", false},
		{"0", "0", false},
		{"f", "f", false},
		// invalid value → falls back to default (false for API_SERVER_ONLY)
		{"invalid", "invalid", false},
		// "yes" was valid in old code — now invalid, falls back to default (false)
		{"yes", "yes", false},
		{"no", "no", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("API_SERVER_ONLY", tt.envValue)
			got := config.APIServerOnly()
			if got != tt.want {
				t.Errorf("APIServerOnly() with API_SERVER_ONLY=%q: got %v, want %v",
					tt.envValue, got, tt.want)
			}
		})
	}
}

// TestGetBool_unset ensures the default value is returned when env var is not set.
func TestGetBool_unset(t *testing.T) {
	orig, wasSet := os.LookupEnv("API_SERVER_ONLY")
	if err := os.Unsetenv("API_SERVER_ONLY"); err != nil {
		t.Fatalf("failed to unset env: %v", err)
	}
	t.Cleanup(func() {
		if wasSet {
			os.Setenv("API_SERVER_ONLY", orig)
		}
	})

	if got := config.APIServerOnly(); got != false {
		t.Errorf("APIServerOnly() with unset env: got %v, want false", got)
	}
}

// TestDataRootDir_pathNormalization ensures the directory path always ends with "/".
func TestDataRootDir_pathNormalization(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{"no trailing slash", "/data/schedules", "/data/schedules/"},
		{"with trailing slash", "/data/schedules/", "/data/schedules/"},
		{"relative path", "./.db", "./.db/"},
		{"already correct", "./mydir/", "./mydir/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DATA_ROOT_DIR", tt.envValue)
			got := config.DataRootDir()
			if got != tt.want {
				t.Errorf("DataRootDir() with DATA_ROOT_DIR=%q: got %q, want %q",
					tt.envValue, got, tt.want)
			}
		})
	}
}

// TestDataRootDir_default ensures the default is "./.db/" (with trailing slash).
func TestDataRootDir_default(t *testing.T) {
	orig, wasSet := os.LookupEnv("DATA_ROOT_DIR")
	if err := os.Unsetenv("DATA_ROOT_DIR"); err != nil {
		t.Fatalf("failed to unset env: %v", err)
	}
	t.Cleanup(func() {
		if wasSet {
			os.Setenv("DATA_ROOT_DIR", orig)
		}
	})

	got := config.DataRootDir()
	if got != "./.db/" {
		t.Errorf("DataRootDir() default: got %q, want %q", got, "./.db/")
	}
}

// TestLogLevel_valid ensures valid log level strings are parsed correctly.
func TestLogLevel_valid(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	level := config.LogLevel()
	if level.String() != "debug" {
		t.Errorf("LogLevel() with LOG_LEVEL=debug: got %q, want %q", level.String(), "debug")
	}
}

// TestLogLevel_invalid ensures an invalid log level falls back to InfoLevel.
func TestLogLevel_invalid(t *testing.T) {
	t.Setenv("LOG_LEVEL", "not-a-level")
	level := config.LogLevel()
	if level.String() != "info" {
		t.Errorf("LogLevel() with invalid LOG_LEVEL: got %q, want %q", level.String(), "info")
	}
}

