include userSettings.env
ifndef VERBOSE
.SILENT:
endif

lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.53.3
	golangci-lint run .
build:
	export DB_USER=${DATABASE_USER} && \
	export DB_PASSWORD=${DATABASE_PASSWORD} && \
	export MQ_USER=${BROKER_USER} && \
	export MQ_PASSWORD=${BROKER_PASSWORD} && \
	docker compose build
test:
	go clean -testcache;
	go test -v -race -count 100 ./configs
	go test -v -race -count 100 ./banner_selector
	go test -v -race -count 100 ./router
integration_test:
	go clean -testcache;
	export DB_USER=${DATABASE_USER} && \
	export DB_PASSWORD=${DATABASE_PASSWORD} && \
	export MQ_USER=${BROKER_USER} && \
	export MQ_PASSWORD=${BROKER_PASSWORD} && \
	docker compose up -d
	# go test -v ./handlers
	go test -v ./database
	go test -v ./message_broker
	go test -v ./integration_tests
	docker compose down --volumes
run:
	export DB_USER=${DATABASE_USER} && \
	export DB_PASSWORD=${DATABASE_PASSWORD} && \
	export MQ_USER=${BROKER_USER} && \
	export MQ_PASSWORD=${BROKER_PASSWORD} && \
	docker compose up -d

down:
	docker compose down --volumes