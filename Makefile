.PHONY: bin
.PHONY: coverage

build:
	mkdir -p bin
	$(shell cd src/ ; go build -a -o app ; cp app ../bin/app ; rm app)

test:
	go test ./...

lint:
	$(shell golangci-lint run)

coverage:
	mkdir -p coverage
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
	go tool cover -func=coverage/coverage.out

upgrade:
	go install github.com/oligot/go-mod-upgrade@latest
	go-mod-upgrade

clean:
	rm -rf bin
	rm -rf coverage

all:
	@make -s clean
	@make -s build
	@make -s lint
	@make -s test
	@make -s coverage
