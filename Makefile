SRC_DIRS=cmd pkg

TESTARGS=-v -race -cover -timeout 20s -cpu 24

VERSION=$(shell git describe || git rev-parse HEAD)
DATE=$(shell date +%s)
LD_FLAGS=-extldflags '-static' -X github.com/stefanoj3/dirstalk/pkg/cmd.Version=$(VERSION) -X github.com/stefanoj3/dirstalk/pkg/cmd.BuildTime=$(DATE)

ifeq ($(CI), true)
TESTARGS=-v -race -coverprofile=coverage.txt -covermode=atomic -timeout 20s -cpu 24
endif

.PHONY: dep
## Fetch dependencies
dep:
	@go mod download

.PHONY: tests
## Execute tests
tests:
	@echo "Executing tests"
	@CGO_ENABLED=1 go test $(TESTARGS) ./...

.PHONY: functional-tests
## Execute functional test
functional-tests: build build-testserver
	./functional-tests.sh

.PHONY: check
## Run checks against the codebase
check:
	docker run -t --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.30-alpine golangci-lint run -v

.PHONY: fix
## Run goimports against the source code
fix:
	docker run --rm -v $(pwd):/data cytopia/goimports -w .

.PHONY: fmt
## Run fmt against the source code
fmt:
	gofmt -s -w $(SRC_DIRS)

.PHONY: release-snapshot
## Creates a release snapshot - requires goreleaser to be available in the $PATH
release-snapshot:
	@echo "Creating release snapshot..."
	@goreleaser --snapshot --rm-dist

.PHONY: release
## Creates a release - requires goreleaser to be available in the $PATH
release:
	@echo "Creating release ..."
	@goreleaser release --skip-publish --rm-dist

.PHONY: build
## Builds binary from source
build:
	CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "$(LD_FLAGS)" -o dist/dirstalk cmd/dirstalk/main.go

.PHONY: build-testserver
## Builds binary for testserver used for the functional tests
build-testserver:
	go build -o dist/testserver cmd/testserver/main.go

.PHONY: out-find
## Search for .out file (profiling) in the repo
out-find:
	@echo "Searching for *.out files"
	find . -name '*.out'

.PHONY: out-delete
## Delete *.out files from the repo
out-delete:
	@echo "Delete *.out files"
	find . -name '*.out' -delete

.PHONY: help
HELP_WIDTH="                       "
## Display makefile help
help:
	@printf "Usage\n";
	@awk '{ \
			if ($$0 ~ /^.PHONY: [a-zA-Z\-\_0-9]+$$/) { \
				helpCommand = substr($$0, index($$0, ":") + 2); \
				if (helpMessage) { \
					printf "  \033[32m%-20s\033[0m %s\n", \
						helpCommand, helpMessage; \
					helpMessage = ""; \
				} \
			} else if ($$0 ~ /^[a-zA-Z\-\_0-9.]+:/) { \
				helpCommand = substr($$0, 0, index($$0, ":")); \
				if (helpMessage) { \
					printf "  \033[32m%-20s\033[0m %s\n", \
						helpCommand, helpMessage; \
					helpMessage = ""; \
				} \
			} else if ($$0 ~ /^##/) { \
				if (helpMessage) { \
					helpMessage = helpMessage"\n"${HELP_WIDTH}substr($$0, 3); \
				} else { \
					helpMessage = substr($$0, 3); \
				} \
			} else { \
				if (helpMessage) { \
					print "\n"${HELP_WIDTH}helpMessage"\n" \
				} \
				helpMessage = ""; \
			} \
		}' \
$(MAKEFILE_LIST)
