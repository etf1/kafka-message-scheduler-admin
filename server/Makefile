VERSION ?= $(shell git describe --always --abbrev=1 --tags --match "v[0-9]*")
LDFLAGS=-ldflags "-X main.version=${VERSION}"
TMPDIR := $(shell mktemp -d)

up:
	docker-compose -p dev up -d kafka scheduler

down:
	docker-compose -p dev down -v

build:
	go build -tags musl -v ./...

.PHONY: bin
bin:
	go build ${LDFLAGS} -tags musl -v -o bin/admin ./cmd/kafka

.PHONY: mini
mini:
	go build ${LDFLAGS} -tags musl -v -o bin/mini ./cmd/mini

start.mini:
	go run ${LDFLAGS} -tags musl -v ./cmd/mini

start:
	go run ${LDFLAGS} -tags musl -v ./cmd/kafka

lint:
	golangci-lint --timeout 5m --build-tags musl run

test:
	go test -v -tags musl -count=1 ./...

test.integration:
	RUN_INTEGRATION_TESTS=yes go test -v -tags musl -count=10 -run=^TestKafkaWatchableStore_AddBuckets ./...

tests: test.integration

# test.integration:
# 	RUN_INTEGRATION_TESTS=yes go test -v -tags musl -count=1 ./...

# tests: lint test test.integration

tests.docker:
	docker-compose -p tests build tests
	BOOTSTRAP_SERVERS="kafka:29092" docker-compose -p tests up tests | tee ${TMPDIR}/tests.result
	docker-compose -p tests down
	bash ./scripts/check_gotest.sh ${TMPDIR}/tests.result
