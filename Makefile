.PHONY: run upx clean clean-db test docker dev-up dev-down web

TARGET = bin/ohana

all: clean test web $(TARGET) upx

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

upx: $(TARGET)
	-upx $(TARGET)

web:
	cd web && \
		yarn && \
		yarn build

run: $(TARGET)
	./$(TARGET)

dev-up: dev-down
	docker run --rm -d \
		--name ohana-postgres-dev \
		-p 127.0.0.1:5432:5432 \
		-e POSTGRES_USER=ohanaAdmin \
		-e POSTGRES_PASSWORD=ohanaMeansFamily \
		-e POSTGRES_DB=ohana \
		postgres:14.2
	docker run --rm -d \
		--name ohana-redis-dev \
		-p 127.0.0.1:6379:6379 \
		redis:7 \
		--requirepass ohanaMeansFamily

dev-down:
	-docker stop ohana-postgres-dev
	-docker stop ohana-redis-dev

dev: dev-up
	go install github.com/codegangsta/gin@latest
	`go env GOPATH`/bin/gin --immediate\
		--port 8000 \
		--appPort 4000 \
		--build cmd/ohana/ \
		--bin ./bin/ohana.gin \
		--buildArgs "-tags osusergo,netgo"

clean:
	rm -rf $(TARGET)
	rm -rf coverage.*

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

docker:
	docker build -f ./.docker/Dockerfile -t ohana .
