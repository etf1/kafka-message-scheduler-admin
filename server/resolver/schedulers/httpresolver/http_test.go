// INTEGRATION TESTS
package httpresolver_test

import (
	"testing"

	"github.com/etf1/kafka-message-scheduler-admin/server/config"
	"github.com/etf1/kafka-message-scheduler-admin/server/helper"
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/httpresolver"
)

func TestDnsResolver(t *testing.T) {
	helper.VerifyIfSkipIntegrationTests(t)

	resolver := httpresolver.NewResolver(config.SchedulersAddr())

	list, err := resolver.List()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(list) != 1 {
		t.Errorf("unexpected list length: %v", len(list))
	}
	for i := 0; i < len(list); i++ {
		sch, ok := list[i].(httpresolver.Scheduler)
		if !ok {
			t.Errorf("unexpected scheduler %v type: %T", sch, sch)
		}
		if v := len(sch.Instances); v != 1 {
			t.Errorf("unexpected instances length for index=%v: %v", i, v)
		}
		if v := len(sch.Instances[0].Topics); v != 1 {
			t.Errorf("unexpected topics length for index=%v: %v", i, v)
		}
	}
}
