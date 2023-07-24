ifndef VERBOSE
.SILENT:
endif

lint:

build:

integration_test:
	go clean -testcache;
	export DB_USER="test" && \
	export DB_PASSWORD="test" && \
	export MQ_USER="test" && \
	export MQ_PASSWORD="test" && \
	docker compose up -d
	go test -v -race -count 100 ./database
	go test -v -race -count 100 ./integration_test
	docker compose down --volumes
run: