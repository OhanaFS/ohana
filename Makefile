TARGET = bin/ohana

UNAME_S := $(shell uname -s)

ifeq ($(UNAME_S),Darwin)
	GOOSNAME := darwin
	LDFLAGS := " \
		-X main.BuildTime=`date -u -I seconds` \
		-X main.GitCommit=`git rev-parse HEAD 2>/dev/null` \
		"
else
	GOOSNAME := linux
	LDFLAGS := " \
		-extldflags=-static \
		-X main.BuildTime=`date -uIs` \
		-X main.GitCommit=`git rev-parse HEAD 2>/dev/null` \
		"
endif


.PHONY: all
all: clean web test $(TARGET) postbuild

$(TARGET): $(shell find . -name '*.go')
	mkdir -p bin
	GOOS=${GOOSNAME} go build -tags osusergo,netgo \
		-ldflags ${LDFLAGS}\
		-o $(TARGET) \
		cmd/ohana/main.go


.PHONY: postbuild
postbuild: $(TARGET)
	strip $(TARGET)

.PHONY: web
web: deps
	cd web && yarn build

.PHONY: clean
clean: deps
	rm -rf $(TARGET) coverage.*
	cd web && yarn clean

.PHONY: deps
deps:
	go mod download -x
	cd web && yarn

.PHONY: run
run: $(TARGET)
	./$(TARGET)

.PHONY: dev-up
dev-up: dev-down
	echo "This command is no longer necessary. Use 'make dev'."
	exit 1

.PHONY: dev-down
dev-down:
	echo "This command is no longer necessary. Use 'make dev'."
	exit 1

.PHONY: dev
dev:
	-rm ./bin/*.gin ./bin/tmp.*
	MY_UID=`id -u` MY_GID=`id -g` \
		docker-compose -f .dev/docker-compose.yaml up --build

.PHONY: format
format:
	go fmt ./...

.PHONY: test
test: web
	go vet ./...
	go test -coverprofile=coverage.out -tags osusergo,netgo ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: docker
docker:
	docker build -f ./.docker/Dockerfile -t ohana .
