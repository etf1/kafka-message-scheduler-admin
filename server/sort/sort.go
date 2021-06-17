package sort

import (
	"sort"
	"strings"

	"github.com/etf1/kafka-message-scheduler/schedule"
)

const (
	Asc Order = iota
	Desc
)

// to reset iota
const (
	Timestamp Field = iota
	ID
	Epoch
)

var (
	DefaultSortBy By = By{
		Timestamp,
		Desc,
	}
)

type Order int

var orderMap = map[string]Order{
	"asc":  Asc,
	"desc": Desc,
}

func (s Order) String() string {
	return [...]string{"asc", "desc"}[s]
}

func ToOrder(order string) Order {
	key := strings.ToLower(order)
	if s, ok := orderMap[key]; ok {
		return s
	}
	return Desc
}

type Field int

var fieldMap = map[string]Field{
	"timestamp": Timestamp,
	"id":        ID,
	"epoch":     Epoch,
}

func (f Field) String() string {
	return [...]string{"timestamp", "id", "epoch"}[f]
}

func ToField(field string) Field {
	key := strings.ToLower(field)
	if s, ok := fieldMap[key]; ok {
		return s
	}
	return Timestamp
}

type By struct {
	Field
	Order
}

func NewSort(arr []schedule.Schedule, sb By) sort.Interface {
	switch sb.Field {
	case Timestamp:
		return ByTimestamp{
			arr,
			sb.Order,
		}
	case ID:
		return ByID{
			arr,
			sb.Order,
		}
	case Epoch:
		return ByEpoch{
			arr,
			sb.Order,
		}
	default:
		return ByTimestamp{
			arr,
			sb.Order,
		}
	}
}

func ToSortBy(s string) By {
	field := ""
	order := ""
	WithoutOrder := 1
	WithOrder := 2

	sortBy := strings.ToLower(strings.TrimSpace(s))

	arr := strings.Split(sortBy, " ")
	if len(arr) == WithoutOrder {
		term := strings.TrimSpace(arr[0])
		if term == "asc" || term == "desc" {
			field = "timestamp"
			order = term
		} else {
			field = term
		}
	} else if len(arr) == WithOrder {
		field = strings.TrimSpace(arr[0])
		order = strings.TrimSpace(arr[1])
	}

	result := By{
		ToField(field),
		ToOrder(order),
	}

	return result
}

type ByTimestamp struct {
	data []schedule.Schedule
	Order
}

func (s ByTimestamp) Len() int { return len(s.data) }
func (s ByTimestamp) Less(i, j int) bool {
	// if equal order by ID asc
	if s.data[i].Timestamp() == s.data[j].Timestamp() {
		return s.data[i].ID() < s.data[j].ID()
	}
	if s.Order == Desc {
		return s.data[i].Timestamp() > s.data[j].Timestamp()
	}
	return s.data[i].Timestamp() < s.data[j].Timestamp()
}
func (s ByTimestamp) Swap(i, j int) { s.data[i], s.data[j] = s.data[j], s.data[i] }

type ByID struct {
	data []schedule.Schedule
	Order
}

func (s ByID) Len() int { return len(s.data) }
func (s ByID) Less(i, j int) bool {
	// if equal order by timestamp asc
	if s.data[i].ID() == s.data[j].ID() {
		return s.data[i].Timestamp() < s.data[j].Timestamp()
	}
	if s.Order == Desc {
		return s.data[i].ID() > s.data[j].ID()
	}
	return s.data[i].ID() < s.data[j].ID()
}
func (s ByID) Swap(i, j int) { s.data[i], s.data[j] = s.data[j], s.data[i] }

type ByEpoch struct {
	data []schedule.Schedule
	Order
}

func (s ByEpoch) Len() int { return len(s.data) }
func (s ByEpoch) Less(i, j int) bool {
	// if equal order by timestamp asc
	if s.data[i].Epoch() == s.data[j].Epoch() {
		return s.data[i].ID() < s.data[j].ID()
	}
	if s.Order == Desc {
		return s.data[i].Epoch() > s.data[j].Epoch()
	}
	return s.data[i].Epoch() < s.data[j].Epoch()
}
func (s ByEpoch) Swap(i, j int) { s.data[i], s.data[j] = s.data[j], s.data[i] }
