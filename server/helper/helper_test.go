package helper_test

import (
	"testing"

	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
)

// TestBleveEscapeTerm verifies that special Bleve query characters are properly escaped.
// This is critical for search safety: user-supplied schedule IDs containing special chars
// (e.g. "invoice+123", "topic/sub") must not be interpreted as Bleve query operators.
func TestBleveEscapeTerm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"plain alphanumeric", "videoID123", "videoID123"},
		// hyphen is a Bleve special char
		{"hyphen", "invoice-1", `invoice\-1`},
		// plus sign (Bleve logical operator)
		{"plus", "a+b", `a\+b`},
		// asterisk (Bleve wildcard)
		{"asterisk", "invoice*", `invoice\*`},
		// slash (common in topic names)
		{"slash", "topic/sub", `topic\/sub`},
		// colon (field:value Bleve syntax)
		{"colon", "field:value", `field\:value`},
		// multiple special chars
		{"multi special", "a+b*c", `a\+b\*c`},
		// space is a special char — internal spaces are escaped, leading/trailing are trimmed
		{"leading trailing spaces", "  hello  ", "hello"},
		{"internal space", "hello world", `hello\ world`},
		// question mark (single-char wildcard in Bleve)
		{"question mark", "sch?dule", `sch\?dule`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helper.BleveEscapeTerm(tt.input)
			if got != tt.want {
				t.Errorf("BleveEscapeTerm(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestSplitTrim verifies comma-separated string splitting with whitespace trimming.
// Used by config.getStrings() for SCHEDULERS_ADDR parsing.
func TestSplitTrim(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"single element", "localhost:8000", []string{"localhost:8000"}},
		{"two elements no spaces", "host1:8000,host2:8000", []string{"host1:8000", "host2:8000"}},
		{"two elements with spaces", "host1:8000, host2:8000", []string{"host1:8000", "host2:8000"}},
		{"spaces around commas", " host1 , host2 , host3 ", []string{"host1", "host2", "host3"}},
		// empty string → slice with one empty element (strings.Split behavior)
		{"empty string", "", []string{""}},
		// empty element between commas
		{"empty between commas", "a,,b", []string{"a", "", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helper.SplitTrim(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("SplitTrim(%q) len = %d, want %d (got %v)", tt.input, len(got), len(tt.want), got)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("SplitTrim(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

