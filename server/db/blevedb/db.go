package blevedb

import (
	"fmt"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/etf1/kafka-message-scheduler-admin/server/db"
	"github.com/etf1/kafka-message-scheduler-admin/server/sort"
	"github.com/etf1/kafka-message-scheduler-admin/server/store"
	"github.com/etf1/kafka-message-scheduler/schedule"
	log "github.com/sirupsen/logrus"
)

const (
	MaxFound   = 1000
	DefaultMax = 300
	ChanSize   = 1000
)

type DB struct {
	store.BatchableStore
	sourceStore store.Watchable
	idxr        indexer
	updtr       updater
}

type Config struct {
	InternalStore store.BatchableStore
	SourceStore   store.Watchable
	Path          string
}

func NewDB(cfg Config) (DB, error) {
	idxr, err := newIndexer(cfg.Path)
	if err != nil {
		return DB{}, err
	}
	go idxr.start()

	updtr := newUpdater(cfg.InternalStore)
	go updtr.start()

	d := DB{
		cfg.InternalStore,
		cfg.SourceStore,
		idxr,
		updtr,
	}

	go d.watch()

	return d, nil
}

func toDocument(sch schedule.Schedule) document {
	s, ok := sch.(store.Schedule)
	if !ok {
		log.Errorf("unexpected type: %T", sch)
		return document{}
	}

	// TODO: manage kafka.Schedule
	return document{
		ID:        s.ID(),
		SortID:    s.ID(),
		Scheduler: s.SchedulerName,
		Epoch:     s.Epoch(),
		Timestamp: s.Timestamp(),
	}
}

// create an id for bleve, we are using a composite key with scheduler name
// in case of there will be a same schedule id for two different scheduler
func bleveID(sch schedule.Schedule) string {
	switch s := sch.(type) {
	case store.Schedule:
		return s.SchedulerName + "|" + s.ID()
	default:
		return sch.ID()
	}
}

func (d DB) upsert(sch store.Schedule) error {
	// upsert in bbolt store
	err := d.updtr.upsert(sch.ID(), sch)
	if err != nil {
		return err
	}

	// upsert in bleve index
	err = d.idxr.upsert(bleveID(sch), toDocument(sch))
	if err != nil {
		return err
	}
	return nil
}

func (d DB) delete(sch schedule.Schedule) error {
	// delete in bbolt store
	err := d.updtr.delete(sch.ID(), sch)
	if err != nil {
		return err
	}

	// delete in bleve index
	err = d.idxr.delete(bleveID(sch))
	if err != nil {
		return err
	}
	return nil
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
			err := d.upsert(evt.Schedule)
			if err != nil {
				log.Errorf("cannot upsert event %+v: %v", evt.Schedule, err)
			}
		case store.DeletedType:
			log.Printf("received delete watch event from store: %+v", evt)
			err := d.delete(evt.Schedule)
			if err != nil {
				log.Errorf("cannot delete event %+v: %v", evt.Schedule, err)
			}
		// the store has been reset need to delete all data
		case store.StoreResetType:
			schedulerName := evt.Schedule.SchedulerName
			log.Printf("received store reset watch event from store for %v: %+v", schedulerName, evt)
			// search all schedules for the reseted scheduler
			_, list, err := d.Search(db.SearchQuery{
				Filter: db.Filter{
					SchedulerName: schedulerName,
				},
				Limit: db.Limit{
					Max: -1,
				},
			})
			if err != nil {
				log.Errorf("cannot find all schedules for %v: %v", schedulerName, err)
				break
			}
			// delete all schedules for the specified scheduler
			for sch := range list {
				err = d.delete(sch)
				if err != nil {
					log.Errorf("reset: cannot delete event %v: %v", evt.Schedule, err)
				}
			}
		}
	}
}

func (d DB) Close() {
	d.idxr.close()
}

func toBleveSort(s sort.By) []string {
	result := []string{}

	sortField := s.Field.String()
	if s.Field == sort.ID {
		sortField = "sort-id"
	}
	if s.Order == sort.Desc {
		sortField = "-" + sortField
	}

	result = append(result, sortField)

	// add secondary sort for equal value, ie: "-timestamp", "id"
	if s.Field != sort.ID {
		result = append(result, "sort-id")
	} else {
		result = append(result, "-"+sort.Timestamp.String())
	}

	return result
}

func toBleveStringQuery(q db.SearchQuery) string {
	// extract logical operator and term "+video" => "+", "video"
	extractTerm := func(term string) (string, string) {
		// default logical operator: + (should)
		lop := "+"
		if len(term) == 1 {
			return lop, term
		}
		if strings.HasPrefix(term, "-") || strings.HasPrefix(term, "+") {
			return term[0:1], term[1:]
		}

		return lop, term
	}

	appendQuery := func(query string, field string, quote bool, s string) string {
		terms := []string{}
		for _, st := range strings.Fields(s) {
			lop, term := extractTerm(st)
			if quote && !strings.Contains(term, "*") {
				term = fmt.Sprintf("%q", term)
			}

			terms = append(terms, lop+field+":"+term)
		}

		if query == "" {
			return strings.Join(terms, " ")
		}

		return query + " " + strings.Join(terms, " ")
	}

	squery := ""

	if q.Filter.ScheduleID != "" {
		squery = appendQuery(squery, "id", true, q.Filter.ScheduleID)
	}

	if q.Filter.SchedulerName != "" {
		squery = appendQuery(squery, "scheduler", true, "+"+q.Filter.SchedulerName)
	}

	if q.Filter.EpochRange.From != 0 {
		squery = appendQuery(squery, "epoch", false, "+"+fmt.Sprintf(">=%v", q.Filter.EpochRange.From))
	}

	if q.Filter.EpochRange.To != 0 {
		squery = appendQuery(squery, "epoch", false, "+"+fmt.Sprintf("<=%v", q.Filter.EpochRange.To))
	}

	return squery
}

func (d DB) Search(q db.SearchQuery) (total int, result chan schedule.Schedule, err error) {
	result = make(chan schedule.Schedule, ChanSize)

	// projected fields
	fields := []string{"scheduler", "id"}
	max := DefaultMax
	if q.Limit.Max == -1 {
		count, err2 := d.idxr.DocCount()
		if err2 == nil {
			max = int(count)
		}
	} else if q.Limit.Max > 0 && q.Limit.Max < MaxFound {
		max = q.Limit.Max
	}

	sortBy := toBleveSort(q.SortBy)
	queryString := toBleveStringQuery(q)

	// bleve search
	var searchQuery query.Query
	if queryString == "" {
		searchQuery = bleve.NewMatchAllQuery()
	} else {
		searchQuery = bleve.NewQueryStringQuery(queryString)
	}

	search := bleve.NewSearchRequest(searchQuery)
	search.SortBy(sortBy)
	search.Size = max
	search.Fields = fields

	log.Warnf("search query='%v' max=%v sort=%v", queryString, max, sortBy)
	start := time.Now()
	searchResults, err := d.idxr.Search(search)
	if err != nil {
		return 0, nil, err
	}

	docCount, err := d.idxr.DocCount()
	fmt.Printf("doc count: %v %v\n", docCount, err)

	fmt.Printf("search done query=%v elapsed=%v: %v\n", searchQuery, time.Since(start), searchResults)

	hitsCount := searchResults.Total

	go func() {
		defer close(result)
		globalStart := time.Now()
		for _, hit := range searchResults.Hits {
			scheduler, ok := hit.Fields["scheduler"].(string)
			if !ok {
				log.Errorf("unexpected scheduler value/type %+v %+v %T", *hit, hit.Fields, hit.Fields["scheduler"])
				continue
			}
			scheduleID, ok := hit.Fields["id"].(string)
			if !ok {
				log.Errorf("unexpected schedule id value/type %+v %+v %T", *hit, hit.Fields, hit.Fields["id"])
				continue
			}
			start := time.Now()
			// get complete schedule object from internal store
			schs, err := d.Get(scheduler, scheduleID)
			log.Warnf("store get one hit %v: %v", scheduleID, time.Until(start))
			if err != nil {
				log.Errorf("unexpected error: %v", err)
				continue
			}
			if len(schs) == 0 {
				log.Errorf("unexpected empty result for %v", scheduleID)
				continue
			}
			result <- schs[0]
		}
		log.Warnf("store get all hits: %v", time.Until(globalStart))
	}()

	return int(hitsCount), result, nil
}
