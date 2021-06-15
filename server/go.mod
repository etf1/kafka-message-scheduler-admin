module github.com/etf1/kafka-message-scheduler-admin/server

go 1.15

require (
	github.com/blevesearch/bleve/v2 v2.0.3
	github.com/blevesearch/bleve_index_api v1.0.0
	github.com/confluentinc/confluent-kafka-go v1.5.2
	github.com/etf1/kafka-message-scheduler v0.0.4-0.20210615142246-56c1d6186d8f // indirect
	github.com/gorilla/mux v1.8.0
	github.com/influxdata/influxdb v1.8.5 // indirect
	github.com/prometheus/common v0.14.0
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.8.1
	github.com/tevjef/go-runtime-metrics v0.0.0-20170326170900-527a54029307
	go.etcd.io/bbolt v1.3.5
	golang.org/x/sys v0.0.0-20210304124612-50617c2ba197 // indirect
)

//replace github.com/etf1/kafka-message-scheduler => /Users/fkarakas/go/src/github.com/etf1/kafka-message-scheduler
