package sort

import (
	"sort"
	"strings"

	"github.com/etf1/kafka-message-scheduler/schedule"
)

const (
	Asc SortOrder = iota
	Desc
)

// to reset iota
const (
	Timestamp SortField = iota
	ID
	Epoch
)

var (
	DefaultSortBy SortBy = SortBy{
		Timestamp,
		Desc,
	}
)

type SortOrder int

var sortOrderMap = map[string]SortOrder{
	"asc":  Asc,
	"desc": Desc,
}

func (s SortOrder) String() string {
	return [...]string{"asc", "desc"}[s]
}

func ToSortOrder(order string) SortOrder {
	key := strings.ToLower(order)
	if s, ok := sortOrderMap[key]; ok {
		return s
	}
	return Desc
}

type SortField int

var sortFieldMap = map[string]SortField{
	"timestamp": Timestamp,
	"id":        ID,
	"epoch":     Epoch,
}

func (s SortField) String() string {
	return [...]string{"timestamp", "id", "epoch"}[s]
}

func ToSortField(field string) SortField {
	key := strings.ToLower(field)
	if s, ok := sortFieldMap[key]; ok {
		return s
	}
	return Timestamp
}

type SortBy struct {
	SortField
	SortOrder
}

func NewSort(arr []schedule.Schedule, sb SortBy) sort.Interface {
	switch sb.SortField {
	case Timestamp:
		return SortByTimestamp{
			arr,
			sb.SortOrder,
		}
	case ID:
		return SortByID{
			arr,
			sb.SortOrder,
		}
	case Epoch:
		return SortByEpoch{
			arr,
			sb.SortOrder,
		}
	default:
		return SortByTimestamp{
			arr,
			sb.SortOrder,
		}
	}
}

func ToSortBy(s string) SortBy {
	field := ""
	order := ""

	sortBy := strings.ToLower(strings.TrimSpace(s))

	arr := strings.Split(sortBy, " ")
	if len(arr) == 1 {
		term := strings.TrimSpace(arr[0])
		if term == "asc" || term == "desc" {
			field = "timestamp"
			order = term
		} else {
			field = term
		}
	} else if len(arr) == 2 {
		field = strings.TrimSpace(arr[0])
		order = strings.TrimSpace(arr[1])
	}

	result := SortBy{
		ToSortField(field),
		ToSortOrder(order),
	}

	return result
}

type SortByTimestamp struct {
	data []schedule.Schedule
	SortOrder
}

func (s SortByTimestamp) Len() int { return len(s.data) }
func (s SortByTimestamp) Less(i, j int) bool {
	// if equal order by ID asc
	if s.data[i].Timestamp() == s.data[j].Timestamp() {
		return s.data[i].ID() < s.data[j].ID()
	}
	if s.SortOrder == Desc {
		return s.data[i].Timestamp() > s.data[j].Timestamp()
	}
	return s.data[i].Timestamp() < s.data[j].Timestamp()
}
func (s SortByTimestamp) Swap(i, j int) { s.data[i], s.data[j] = s.data[j], s.data[i] }

type SortByID struct {
	data []schedule.Schedule
	SortOrder
}

func (s SortByID) Len() int { return len(s.data) }
func (s SortByID) Less(i, j int) bool {
	// if equal order by timestamp asc
	if s.data[i].ID() == s.data[j].ID() {
		return s.data[i].Timestamp() < s.data[j].Timestamp()
	}
	if s.SortOrder == Desc {
		return s.data[i].ID() > s.data[j].ID()
	}
	return s.data[i].ID() < s.data[j].ID()
}
func (s SortByID) Swap(i, j int) { s.data[i], s.data[j] = s.data[j], s.data[i] }

type SortByEpoch struct {
	data []schedule.Schedule
	SortOrder
}

func (s SortByEpoch) Len() int { return len(s.data) }
func (s SortByEpoch) Less(i, j int) bool {
	// if equal order by timestamp asc
	if s.data[i].Epoch() == s.data[j].Epoch() {
		return s.data[i].ID() < s.data[j].ID()
	}
	if s.SortOrder == Desc {
		return s.data[i].Epoch() > s.data[j].Epoch()
	}
	return s.data[i].Epoch() < s.data[j].Epoch()
}
func (s SortByEpoch) Swap(i, j int) { s.data[i], s.data[j] = s.data[j], s.data[i] }
