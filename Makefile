#!make
.DEFAULT_GOAL := dc-test

# Makefile target args
args = $(filter-out $@,$(MAKECMDGOALS))

# docker cmds
dc-up:
	docker compose up -d app

dc-build:
	docker compose build --no-cache app

dc-stop:
	docker compose stop app

dc-down:
	docker compose down

dc-delete-volumes:
	docker compose down --volumes

dc-logs:
	docker compose logs -f $(args)

dc-test: dc-up
	docker compose exec app go test -count=1 -v ./...


dc-test-cov: dc-up
	docker compose exec app go test -coverprofile=artifacts/coverage.out ./...
	docker compose exec app go tool cover -html=artifacts/coverage.out -o artifacts/coverage.html 