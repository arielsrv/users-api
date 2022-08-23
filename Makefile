.PHONY: bin
.PHONY: coverage

build:
	/bin/bash -c \
	"(mkdir -p bin ; cd src/ ; go build -o app ; mv app ../bin)"

test:
	go test ./...

lint:
	$(shell golangci-lint run)

gocritic:
	gocritic check ./...

coverage:
	mkdir -p coverage
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
	go tool cover -func=coverage/coverage.out

upgrade:
	go mod tidy
	go install github.com/oligot/go-mod-upgrade@latest
	go-mod-upgrade

clean:
	rm -rf bin/
	rm -rf coverage/

all:
	@make -s clean
	@make -s build
	@make -s lint
	@make -s test
	@make -s coverage
