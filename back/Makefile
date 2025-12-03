BINARY_NAME=fantasy-back

GO=go
GO_BUILD=$(GO) build
GO_RUN=$(GO) run
GO_TEST=$(GO) test

all: build

build:
	$(GO_BUILD) -o $(BINARY_NAME) ./cmd/main.go

run:
	$(GO_RUN) ./cmd/main.go

test:
	$(GO_TEST) ./...

clean:
	rm -f $(BINARY_NAME)

run_build: build
	./$(BINARY_NAME)