.PHONY: run upx clean clean-db test docker dev-up dev-down web

TARGET = bin/ohana

all: clean web test $(TARGET) upx

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
	docker-compose -f .dev/docker-compose.yaml up -d

dev-down:
	-docker-compose -f .dev/docker-compose.yaml down

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
	cd web && yarn && yarn clean

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

docker:
	docker build -f ./.docker/Dockerfile -t ohana .
