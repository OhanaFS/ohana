TARGET = bin/ohana

.PHONY: all
all: clean web test $(TARGET) postbuild

$(TARGET): $(shell find . -name '*.go')
	mkdir -p bin
	GOOS=linux go build -tags osusergo,netgo \
		-ldflags " \
			-extldflags=-static \
			-X main.BuildTime=`date -uIs` \
			-X main.GitCommit=`git rev-parse HEAD 2>/dev/null` \
		" \
		-o $(TARGET) \
		cmd/ohana/main.go

.PHONY: postbuild
postbuild: $(TARGET)
	strip $(TARGET)
	upx $(TARGET)

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
	docker-compose -f .dev/docker-compose.yaml up -d

.PHONY: dev-down
dev-down:
	-docker-compose -f .dev/docker-compose.yaml down

.PHONY: dev
dev: dev-up
	go install github.com/codegangsta/gin@latest
	`go env GOPATH`/bin/gin \
		--immediate \
		--port 8000 \
		--appPort 4000 \
		--build cmd/ohana/ \
		--bin ./bin/ohana.gin \
		--buildArgs "-tags osusergo,netgo"

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
