package sort_test

import (
	stdsort "sort"
	"testing"
	"time"

	sorter "github.com/etf1/kafka-message-scheduler-admin/server/sort"
	"github.com/etf1/kafka-message-scheduler/schedule"
	simple_sched "github.com/etf1/kafka-message-scheduler/schedule/simple"
)

// TestToSortBy validates all parsing branches: empty, order-only, field-only, field+order,
// unknown values, and case-insensitivity.
func TestToSortBy(t *testing.T) {
	tests := []struct {
		input         string
		expectedField sorter.Field
		expectedOrder sorter.Order
	}{
		// empty → default (timestamp desc)
		{"", sorter.Timestamp, sorter.Desc},
		// order-only shortcuts
		{"asc", sorter.Timestamp, sorter.Asc},
		{"desc", sorter.Timestamp, sorter.Desc},
		// field-only (order defaults to Desc)
		{"timestamp", sorter.Timestamp, sorter.Desc},
		{"id", sorter.ID, sorter.Desc},
		{"epoch", sorter.Epoch, sorter.Desc},
		// field + order
		{"timestamp asc", sorter.Timestamp, sorter.Asc},
		{"timestamp desc", sorter.Timestamp, sorter.Desc},
		{"id asc", sorter.ID, sorter.Asc},
		{"id desc", sorter.ID, sorter.Desc},
		{"epoch asc", sorter.Epoch, sorter.Asc},
		{"epoch desc", sorter.Epoch, sorter.Desc},
		// unknown field → fallback to Timestamp
		{"unknown", sorter.Timestamp, sorter.Desc},
		// unknown order → fallback to Desc
		{"timestamp invalid", sorter.Timestamp, sorter.Desc},
		// case-insensitive
		{"TIMESTAMP ASC", sorter.Timestamp, sorter.Asc},
		{"ID DESC", sorter.ID, sorter.Desc},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sorter.ToSortBy(tt.input)
			if got.Field != tt.expectedField {
				t.Errorf("ToSortBy(%q).Field = %v, want %v", tt.input, got.Field, tt.expectedField)
			}
			if got.Order != tt.expectedOrder {
				t.Errorf("ToSortBy(%q).Order = %v, want %v", tt.input, got.Order, tt.expectedOrder)
			}
		})
	}
}

// TestSort_tieBreakerTimestamp: when two schedules share the same Timestamp,
// the tiebreaker must always be ID ascending — regardless of the requested sort order.
func TestSort_tieBreakerTimestamp(t *testing.T) {
	sameTime := time.Unix(1_000_000, 0)

	// "sch-a" < "sch-b" lexicographically → after tiebreak, a comes first.
	schA := simple_sched.NewSchedule("sch-a", sameTime, sameTime)
	schB := simple_sched.NewSchedule("sch-b", sameTime, sameTime)

	for _, order := range []sorter.Order{sorter.Asc, sorter.Desc} {
		t.Run("timestamp/order="+order.String(), func(t *testing.T) {
			// Start with b before a.
			arr := []schedule.Schedule{schB, schA}
			stdsort.Sort(sorter.NewSort(arr, sorter.By{Field: sorter.Timestamp, Order: order}))

			// Tiebreaker: ID asc, regardless of sort direction.
			if arr[0].ID() != "sch-a" || arr[1].ID() != "sch-b" {
				t.Errorf("timestamp tiebreak (order=%v): got [%s, %s], want [sch-a, sch-b]",
					order, arr[0].ID(), arr[1].ID())
			}
		})
	}
}

// TestSort_tieBreakerID: when two schedules share the same ID,
// the tiebreaker is Timestamp ascending — regardless of the requested sort order.
func TestSort_tieBreakerID(t *testing.T) {
	t1 := time.Unix(1_000_000, 0)
	t2 := time.Unix(2_000_000, 0) // t2 > t1

	// Both have the same ID "sch-1" but different timestamps.
	schOld := simple_sched.NewSchedule("sch-1", t1, t1)
	schNew := simple_sched.NewSchedule("sch-1", t2, t2)

	for _, order := range []sorter.Order{sorter.Asc, sorter.Desc} {
		t.Run("id/order="+order.String(), func(t *testing.T) {
			// Start with newer before older.
			arr := []schedule.Schedule{schNew, schOld}
			stdsort.Sort(sorter.NewSort(arr, sorter.By{Field: sorter.ID, Order: order}))

			// Tiebreaker: Timestamp asc → older first (t1.Unix() < t2.Unix()).
			if arr[0].Timestamp() != t1.Unix() || arr[1].Timestamp() != t2.Unix() {
				t.Errorf("id tiebreak (order=%v): got timestamps [%v, %v], want [%v, %v]",
					order, arr[0].Timestamp(), arr[1].Timestamp(), t1.Unix(), t2.Unix())
			}
		})
	}
}

// TestSort_tieBreakerEpoch: when two schedules share the same Epoch,
// the tiebreaker is ID ascending — regardless of the requested sort order.
func TestSort_tieBreakerEpoch(t *testing.T) {
	sameEpoch := time.Unix(1_000_000, 0)
	t1 := time.Unix(1_000_001, 0)
	t2 := time.Unix(1_000_002, 0)

	// Same epoch, different IDs and timestamps.
	schA := simple_sched.NewSchedule("sch-a", sameEpoch, t1)
	schB := simple_sched.NewSchedule("sch-b", sameEpoch, t2)

	for _, order := range []sorter.Order{sorter.Asc, sorter.Desc} {
		t.Run("epoch/order="+order.String(), func(t *testing.T) {
			// Start with b before a.
			arr := []schedule.Schedule{schB, schA}
			stdsort.Sort(sorter.NewSort(arr, sorter.By{Field: sorter.Epoch, Order: order}))

			// Tiebreaker: ID asc → a before b.
			if arr[0].ID() != "sch-a" || arr[1].ID() != "sch-b" {
				t.Errorf("epoch tiebreak (order=%v): got [%s, %s], want [sch-a, sch-b]",
					order, arr[0].ID(), arr[1].ID())
			}
		})
	}
}

