package helper

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	lorem "github.com/drhodes/golorem"
	log "github.com/sirupsen/logrus"
)

const (
	MaxRand  = 1000000
	LoremMin = 50
	LoremMax = 100
)

func Lipsum() string {
	return lorem.Paragraph(LoremMin, LoremMax)
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
	return fmt.Sprintf("%s%v", prefix, RandNum())
}

func RandNum() int64 {
	return RandNumWithMax(MaxRand)
}

func IfElse(cond bool, then, otherwise int) int {
	if cond {
		return then
	}
	return otherwise
}

func RandNumWithMax(max int64) int64 {
	bg := big.NewInt(max)

	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}

	return n.Int64()
}

func LogErr(err error) {
	if err != nil {
		log.Errorf("an error has occurred: %v", err)
	}
}

func BleveEscapeTerm(term string) string {
	v := strings.TrimSpace(term)
	// special chars to be escaped for bleve
	specialChars := "+-=&|><!(){}[]^\"~*?:\\/ "
	result := ""
	for _, c := range v {
		if strings.ContainsRune(specialChars, c) {
			result += "\\"
		}
		result += string(c)
	}
	return result
}
