package kafka

import (
	"github.com/etf1/kafka-message-scheduler-admin/server/resolver/schedulers/httpresolver"
	"github.com/etf1/kafka-message-scheduler-admin/server/store/kafka"
	log "github.com/sirupsen/logrus"
)

func SchedulesTopics(s httpresolver.Scheduler) []string {
	return s.Topics()
}

func HistoryTopic(s httpresolver.Scheduler) []string {
	return []string{s.History()}
}

func NewWatchableStoreFromResolver(r httpresolver.Resolver, topics func(s httpresolver.Scheduler) []string) (kafka.WatchableStore, error) {
	lst, err := r.List()
	if err != nil {
		return kafka.WatchableStore{}, err
	}

	buckets := make([]kafka.Bucket, 0)

	for _, sch := range lst {
		s, ok := sch.(httpresolver.Scheduler)
		if !ok {
			log.Errorf("unable to cast scheduler: %T", sch)
			continue
		}
		buckets = append(buckets, kafka.Bucket{
			Name:             s.Name(),
			BootstrapServers: s.BootstrapServers(),
			Topics:           topics(s),
		})
	}

	return kafka.NewWatchableStore(buckets)
}
