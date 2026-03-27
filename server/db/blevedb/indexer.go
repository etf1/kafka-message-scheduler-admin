package blevedb

import (
	"fmt"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/simple"
	log "github.com/sirupsen/logrus"
)

const (
	batchSize = 1000
)

type eventType int

const (
	upsertType eventType = iota
	deleteType
)

type document struct {
	ID          string `json:"id"`
	SortID      string `json:"sort-id"`
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

func newIndexer(path string) (*indexer, error) {
	// Try to open an existing index first (server restart case).
	index, err := bleve.Open(path)
	if err == nil {
		log.Infof("opened existing bleve index at %s", path)
		return &indexer{
			make(chan event, MaxChanSize),
			index,
		}, nil
	}

	// Index does not exist yet — create a new one.
	log.Infof("creating new bleve index at %s", path)

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	// a generic reusable mapping for simple text
	simpleFieldMapping := bleve.NewTextFieldMapping()
	simpleFieldMapping.Analyzer = simple.Name

	// mapping
	mapping := bleve.NewIndexMapping()
	mapping.DefaultMapping = bleve.NewDocumentMapping()
	mapping.DefaultMapping.AddFieldMappingsAt("id", simpleFieldMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("scheduler", keywordFieldMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("sort-id", keywordFieldMapping)
	mapping.DefaultMapping.AddFieldMappingsAt("epoch", bleve.NewNumericFieldMapping())
	mapping.DefaultMapping.AddFieldMappingsAt("timestamp", bleve.NewNumericFieldMapping())

	index, err = bleve.New(path, mapping)
	if err != nil {
		return nil, err
	}

	return &indexer{
		make(chan event, MaxChanSize),
		index,
	}, nil
}

func (i *indexer) close() {
	close(i.input)
	i.input = nil
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

			toDocument := func(data interface{}) (document, error) {
				if data == nil {
					return document{}, fmt.Errorf("nil object")
				}
				doc, ok := data.(document)
				if !ok {
					return document{}, fmt.Errorf("unexpected object type: %T", data)
				}
				return doc, nil
			}

			switch evt.eventType {
			case upsertType:
				doc, err := toDocument(evt.data)
				if err != nil {
					log.Error(err)
					break
				}
				log.Printf("batch index: %+v", doc)
				err = batch.Index(evt.id, doc)
				if err != nil {
					log.Errorf("index batch failed: %v", err)
					break
				}
			case deleteType:
				log.Printf("batch delete with id: %v", evt.id)
				batch.Delete(evt.id)
			}
			counter++
			if counter%batchSize == 0 {
				indexBatch()
				log.Warnf("indexed %v documents", counter)
			}
		case <-timeout.C:
			log.Tracef("input channel timeout")
			if batch.Size() != 0 {
				indexBatch()
				log.Debugf("indexed %v documents", counter)
			}
		}
	}
}

func (i indexer) upsert(id string, data document) {
	if i.input == nil {
		return
	}
	i.input <- event{
		upsertType,
		id,
		data,
	}
}

func (i indexer) delete(id string) {
	if i.input == nil {
		return
	}
	i.input <- event{
		eventType: deleteType,
		id:        id,
	}
}

