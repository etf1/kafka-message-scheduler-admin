package blevedb

import (
	"fmt"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	index "github.com/blevesearch/bleve_index_api"
	log "github.com/sirupsen/logrus"
)

const (
	batch_size = 1000
)

type eventType int

const (
	upsertType eventType = iota
	deleteType
)

type document struct {
	Id          string `json:"id"`
	Scheduler   string `json:"scheduler"`
	Epoch       int64  `json:"epoch"`
	Timestamp   int64  `json:"timestamp"`
	Topic       string `json:"topic"`
	TargetTopic string `json:"target-topic"`
	TargetKey   string `json:"target-key"`
}

type event struct {
	eventType
	id   string
	data interface{}
}

type indexer struct {
	input chan event
	bleve.Index
}

func newIndexer(path string) (indexer, error) {
	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	// mapping
	mapping := bleve.NewIndexMapping()
	mapping.DefaultMapping = bleve.NewDocumentMapping()
	mapping.DefaultMapping.AddFieldMappingsAt("id", keywordFieldMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("scheduler", keywordFieldMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("epoch", bleve.NewNumericFieldMapping())
	mapping.DefaultMapping.AddFieldMappingsAt("timestamp", bleve.NewNumericFieldMapping())

	/*
		indexMapping := bleve.NewIndexMapping()

		storeScheduleMapping := bleve.NewDocumentMapping()
		storeScheduleMapping.AddFieldMappingsAt("scheduler", keywordFieldMapping)

		scheduleMapping := bleve.NewDocumentMapping()
		scheduleMapping.AddFieldMappingsAt("id", keywordFieldMapping)
		storeScheduleMapping.AddSubDocumentMapping("schedule", scheduleMapping)

		//indexMapping.AddDocumentMapping("root", storeScheduleMapping)

		indexMapping.DefaultMapping = storeScheduleMapping
	*/

	index, err := bleve.New(path, mapping)
	if err != nil {
		return indexer{}, err
	}

	return indexer{
		make(chan event, 10000),
		index,
	}, nil
}

func (i indexer) close() {
	close(i.input)
}

func (i indexer) start() {
	defer log.Printf("indexer closed")

	duration := 500 * time.Millisecond
	timeout := time.NewTimer(duration)
	defer timeout.Stop()

	counter := 0
	batch := i.NewBatch()

	indexBatch := func() {
		log.Printf("batch indexing %v documents", counter)
		err := i.Batch(batch)
		if err != nil {
			log.Printf("batch indexing failed : %v", err)
		}
		//batch.Reset()
		batch = i.NewBatch()
	}

loop:
	for {
		timeout.Reset(duration)
		select {
		case evt, ok := <-i.input:
			log.Printf("indexer: received event from input channel")
			if !ok {
				log.Printf("input channel closed")
				indexBatch()
				break loop
			}

			doc, ok := evt.data.(document)
			if !ok {
				log.Errorf("unexpected document object: %T", evt.data)
			}

			switch evt.eventType {
			case upsertType:
				log.Printf("batch index: %+v", doc)
				err := batch.Index(evt.id, doc)
				if err != nil {
					log.Errorf("index batch failed: %v", err)
					break
				}
			case deleteType:
				log.Printf("batch delete: %v", doc)
				batch.Delete(evt.id)
			}
			counter++
			if counter%batch_size == 0 {
				indexBatch()
				//batch = i.NewBatch()
				log.Warnf("indexed %v documents", counter)
			}
		case <-timeout.C:
			log.Printf("input channel timeout")
			if batch.Size() != 0 {
				indexBatch()
				log.Warnf("indexed %v documents", counter)
			}
		}
	}
}

func (i indexer) upsertDocument(id string, data document) error {
	if i.input == nil {
		return fmt.Errorf("indexer not initialized or closed")
	}
	i.input <- event{
		upsertType,
		id,
		data,
	}
	return nil
}

func (i indexer) deleteDocument(id string, data document) error {
	if i.input == nil {
		return fmt.Errorf("indexer not initialized or closed")
	}
	i.input <- event{
		deleteType,
		id,
		data,
	}
	return nil
}

func describeDocument(doc index.Document) string {
	rv := struct {
		ID     string                 `json:"id"`
		Fields map[string]interface{} `json:"fields"`
	}{
		ID:     doc.ID(),
		Fields: map[string]interface{}{},
	}
	doc.VisitFields(func(field index.Field) {
		var newval interface{}
		switch field := field.(type) {
		case index.TextField:
			newval = field.Text()
		case index.NumericField:
			n, err := field.Number()
			if err == nil {
				newval = n
			}
		case index.DateTimeField:
			d, err := field.DateTime()
			if err == nil {
				newval = d.Format(time.RFC3339Nano)
			}
		}
		existing, existed := rv.Fields[field.Name()]
		if existed {
			switch existing := existing.(type) {
			case []interface{}:
				rv.Fields[field.Name()] = append(existing, newval)
			case interface{}:
				arr := make([]interface{}, 2)
				arr[0] = existing
				arr[1] = newval
				rv.Fields[field.Name()] = arr
			}
		} else {
			rv.Fields[field.Name()] = newval
		}
	})

	return fmt.Sprintf("doc: %+v\n", rv)
}
