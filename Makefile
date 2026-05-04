GOFILES := $(shell find . -path './.tools' -prune -o -name '*.go' -print)

.PHONY: fmt vet test race arch check

fmt:
	@test -z "$$(gofmt -l $(GOFILES))"

vet:
	@go vet ./...

test:
	@go test ./...

race:
	@go test -race ./...

arch:
	@go test ./internal/arch -run 'ImportBoundaries|CoreIsPure|CoreTransitiveDepsAreStdlib|AppDoesNotImportAdapters|PortsAreUseCaseOwned|PortsAreNarrow|ProviderBoundary|CLIIsThin'

check: fmt vet arch test race
