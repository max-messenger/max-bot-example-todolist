.PHONY: test_unit
test_unit:
	go clean -testcache && go test -v -race -run Unit ./...

.PHONY: test_integration
test_integration:
	go clean -testcache && go test -v -race -run Integration ./...
	
.PHONY: test_cover
test_cover:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint: 
	golangci-lint run

.PHONY: dev_env_up dev_env_down
dev_env_up:
	docker-compose pull
	docker-compose up -d
dev_env_down:
	docker-compose  down -v --remove-orphans

.PHONY: generate
generate: 
	go generate -x ./...

migrate-new: ## Create migration
	sql-migrate new -config=db/conf.yml -env=local $(name)