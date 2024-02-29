
.PHONY: server/start
server/start:
	docker-compose -f docker-compose.yml up -d  --remove-orphans

.PHONY: server/stop
server/stop:
	docker-compose -f docker-compose.yml down

.PHONY: server/restart
server/restart:
	docker-compose -f docker-compose.yml restart restapi db-migration
################
# BUILD BINARY
################

GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0

.PHONY: build/restapi
build/restapi:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o build/restapi cmd/main.go