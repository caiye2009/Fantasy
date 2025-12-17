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

# ES Management Tools

# ES Clear - Delete indexes
# Usage:
#   make esclear              # Delete all indexes
#   make esclear domain=client # Delete specific domain index
esclear:
ifdef domain
	$(GO_RUN) ./cmd/tools/esclear.go $(domain)
else
	$(GO_RUN) ./cmd/tools/esclear.go
endif

# ES Reindex - Rebuild indexes from database
# Usage:
#   make esreindex                    # Reindex all supported domains
#   make esreindex domain=client      # Reindex specific domain
#   make esreindex domain="plan process product supplier" # Reindex multiple domains
# Supported domains: client, material, process, product, supplier, order, plan
esreindex:
ifdef domain
	$(GO_RUN) ./cmd/tools/esreindex.go $(domain)
else
	$(GO_RUN) ./cmd/tools/esreindex.go
endif

# Legacy command (deprecated, use esreindex domain=client instead)
reindex-clients:
	$(GO_RUN) ./cmd/tools/reindex_clients.go