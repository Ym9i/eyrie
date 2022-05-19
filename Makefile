.PHONY: test build

export GOPROXY=https://goproxy.cn,direct
export GO111MODULE=on

test:
	gofmt -w -s .
	go test -count=1 -v ./config/...

	@echo "test done!"
build:
	gofmt -w -s .
	go clean
	go build -v -p 4 -o ./bin/ecs pkg/main.go
	# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -p 4 -o ./bin/ecs pkg/main.go

	@echo "build successfully!"
