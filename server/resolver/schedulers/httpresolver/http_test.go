// INTEGRATION TESTS
package httpresolver_test

import (
	"os"
	"testing"

	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/httpresolver"
)

func getSchedulerHost() string {
	if helper.IsRunningInDocker() {
		return "kafka-message-scheduler"
	}
	if v := os.Getenv("MINI_SCHEDULER_HOST"); v != "" {
		return v
	}
	return "localhost:8000"
}

func TestDnsResolver(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	resolver := httpresolver.NewResolver([]string{getSchedulerHost()})

	list, err := resolver.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("unexpected list length: %v", len(list))
	}
	for i := 0; i < len(list); i++ {
		sch, ok := list[i].(httpresolver.Scheduler)
		if !ok {
			t.Fatalf("unexpected scheduler %v type: %T", sch, sch)
		}
		if v := len(sch.Instances); v != 1 {
			t.Fatalf("unexpected instances length for index=%v: %v", i, v)
		}
		if v := len(sch.Instances[0].Topics); v != 1 {
			t.Fatalf("unexpected topics length for index=%v: %v", i, v)
		}
	}
}
