package decoder

import "github.com/etf1/kafka-message-scheduler/schedule"

type Decoder interface {
	Decode(s schedule.Schedule) (schedule.Schedule, error)
}
