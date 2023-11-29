BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN_CALENDAR) -config ./configs/calendar_config.yaml
	$(BIN_SCHEDULER) -config ./configs/scheduler_config.yaml
	$(BIN_SENDER) -config ./configs/sender_config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN_CALENDAR) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

lint: install-lint-deps
	golangci-lint run ./...

gen-grpc:
	buf generate

up:
	docker-compose -f ./deployments/docker-compose.yaml up -d

down:
	docker-compose -f ./deployments/docker-compose.yml down

integration-tests:
		set -e ;\
    docker-compose -f ./deployments/docker-compose.test.yml up --build -d ;\
    test_status_code=0 ;\
    docker-compose -f ./deployments/docker-compose.test.yml run test_itgr go test -tags integration ./test_integration || test_status_code=$$? ;\
    docker-compose -f ./deployments/docker-compose.test.yml down ;\
    exit $$test_status_code ;\

.PHONY: build run build-img run-img version test lint
