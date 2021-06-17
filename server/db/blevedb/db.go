package blevedb

import (
	"fmt"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/etf1/kafka-message-scheduler-admin/server/db"
	"github.com/etf1/kafka-message-scheduler-admin/server/sort"
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler/schedule"
	log "github.com/sirupsen/logrus"
)

type DB struct {
	store.BatchableStore
	sourceStore store.Watchable
	idx         indexer
	updtr       updater
}

type Config struct {
	InternalStore store.BatchableStore
	SourceStore   store.Watchable
	Path          string
}

func NewDB(cfg Config) (DB, error) {
	idx, err := newIndexer(cfg.Path)
	if err != nil {
		return DB{}, err
	}
	go idx.start()

	updtr := newUpdater(cfg.InternalStore)
	go updtr.start()

	d := DB{
		cfg.InternalStore,
		cfg.SourceStore,
		idx,
		updtr,
	}

	go d.watch()

	return d, nil
}

func toDocument(sch store.Schedule) document {
	// TODO: manage kafka.Schedule
	return document{
		Id:        sch.ID(),
		Scheduler: sch.SchedulerName,
		Epoch:     sch.Epoch(),
		Timestamp: sch.Timestamp(),
	}
}

func (d DB) watch() {
	defer log.Printf("watcher closed")

	watchChan, err := d.sourceStore.Watch()
	if err != nil {
		log.Errorf("cannot get watch channel: %v", err)
	}

	for evt := range watchChan {
		log.Printf("received watch event from store: %+v", evt)

		switch evt.EventType {
		case store.UpsertType:
			log.Printf("received upsert watch event from store: %+v", evt)
			d.updtr.upsertData(evt.Schedule.Schedule.ID(), evt.Schedule)
			d.idx.upsertDocument(evt.ID(), toDocument(evt.Schedule))
		case store.DeletedType:
			log.Printf("received delete watch event from store: %+v", evt)
			d.updtr.deleteData(evt.Schedule.Schedule.ID(), evt.Schedule)
			d.idx.deleteDocument(evt.ID(), toDocument(evt.Schedule))
		}
	}
}

func (d DB) Close() {
	d.idx.close()
}

func toBleveSort(s sort.SortBy) []string {
	result := []string{}

	sortField := s.SortField.String()
	if s.SortOrder == sort.Desc {
		sortField = "-" + sortField
	}

	result = append(result, sortField)

	// add secondary sort for equal value, ie: "-timestamp", "id"
	if s.SortField != sort.ID {
		result = append(result, sort.ID.String())
	} else {
		result = append(result, "-"+sort.Timestamp.String())
	}

	return result
}

func (d DB) Search(q db.SearchQuery) (int, chan schedule.Schedule, error) {
	result := make(chan schedule.Schedule, 1000)

	// projected fields
	fields := []string{"scheduler"}
	squery := ""
	max := 100
	if q.Limit.Max > 0 && q.Limit.Max < 1000 {
		max = q.Limit.Max
	}
	sortBy := toBleveSort(q.SortBy)

	appendQuery := func(query string, field string, s string) string {
		if field != "" {
			field = field + ":"
		}
		if !strings.Contains(s, "*") {
			s = "*" + s + "*"
		}
		if query == "" {
			return field + s
		}
		if field == "" {
			return query + " " + s
		}
		return query + " " + field + s
	}

	if q.Filter.ScheduleID != "" {
		squery = appendQuery(squery, "+id", q.Filter.ScheduleID)
	}

	if q.Filter.SchedulerName != "" {
		squery = appendQuery(squery, "+scheduler", q.Filter.SchedulerName)
	}

	if squery == "" {
		squery = "*"
	}

	// bleve search
	query := bleve.NewQueryStringQuery(squery)
	search := bleve.NewSearchRequest(query)
	search.SortBy(sortBy)
	search.Size = max
	search.Fields = fields
	log.Printf("search query='%v' max=%v sort=%v", squery, max, sortBy)
	start := time.Now()
	searchResults, err := d.idx.Search(search)
	if err != nil {
		return 0, nil, err
	}

	fmt.Printf("search done elapsed=%v : %v\n", time.Since(start), searchResults)

	hitsCount := searchResults.Total

	go func() {
		defer close(result)

		for _, hit := range searchResults.Hits {
			scheduler, ok := hit.Fields["scheduler"].(string)
			if !ok {
				log.Errorf("unexpected scheduler value/type %+v %+v %T", *hit, (*hit).Fields, (*hit).Fields["Scheduler"])
				continue
			}
			schs, err := d.Get(scheduler, hit.ID)
			if err != nil {
				log.Errorf("unexpected error: %v", err)
				continue
			}
			if len(schs) == 0 {
				log.Errorf("unexpected empty result for %v", hit.ID)
				continue
			}
			result <- schs[0]
		}
	}()

	// uint64, err := d.idx.DocCount()
	// if err != nil {
	// 	return 0, nil, err
	// }

	// total := int(uint64)

	// fmt.Printf("@@@ total=%v", total)

	// doc, err := d.idx.Document("schedule-2")
	// if err != nil {
	// 	fmt.Printf("%v", err)
	// }

	// fmt.Printf("@@@ %+v", describeDocument(doc))

	return int(hitsCount), result, nil
}
