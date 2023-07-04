unit-test:
	go test -race ./...

int-test:
	$(MAKE) int-test-db-down
	$(MAKE) int-test-db-up
	go test ./... --tags=integration

int-test-db-up:
	docker-compose -f ./internal/integration_tests/docker-compose.environment.yml up -d --wait

int-test-db-down:
	docker-compose -f ./internal/integration_tests/docker-compose.environment.yml down

local-db-up:
	docker-compose -f ./local/docker-compose.environment.yml up -d --wait

local-db-down:
	docker-compose -f ./local/docker-compose.environment.yml down

run-gateway:
	$(MAKE) local-db-down
	$(MAKE) local-db-up
	go run cmd/gateway/main.go

run-storage:
	go run cmd/gateway/main.go

setup:
	go install github.com/vektra/mockery/v2@v2.20.2

clean/mocks:
	find ./pkg/mocks/* -exec rm -rf {} \; || true

generate/mocks: clean/mocks
	mockery --all --case snake --dir ./pkg --output ./pkg/mocks

