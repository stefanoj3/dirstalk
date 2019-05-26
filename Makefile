SRC_DIRS=cmd pkg

TESTARGS=-v -race -cover

VERSION=$(shell git describe || git rev-parse HEAD)
DATE=$(shell date +%s)
LD_FLAGS=-extldflags '-static' -X github.com/stefanoj3/dirstalk/pkg/cmd.Version=$(VERSION) -X github.com/stefanoj3/dirstalk/pkg/cmd.BuildTime=$(DATE)

ifeq ($(TRAVIS), true)
TESTARGS=-v -race -coverprofile=coverage.txt -covermode=atomic
endif

.PHONY: dep
## Fetch dependencies
dep:
	@go get -u github.com/golang/dep/cmd/dep
	@go get -u golang.org/x/tools/cmd/goimports
	@go get -u golang.org/x/lint/golint
	@go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@go get -u github.com/securego/gosec/cmd/gosec
	@dep ensure

.PHONY: test
## Execute tests
test:
	@echo "Executing tests"
	@CGO_ENABLED=1 go test $(TESTARGS) ./...


.PHONY: check
## Run checks against the codebase
check:
	@golint -set_exit_status .
	@go vet ./...
	@goimports -l $(SRC_DIRS) | tee /dev/tty | xargs -I {} test -z {}
	@golangci-lint run
	@gosec ./...

.PHONY: fix
## Run goimports against the source code
fix:
	@goimports -w $(SRC_DIRS)

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

.PHONY: help
## Display this help screen - requires gawk
help:
	@gawk 'match($$0, /^## (.*)/, a) \
		{ getline x; x = gensub(/(.+:) .+/, "\\1", "g", x) ; \
		printf "\033[36m%-30s\033[0m %s\n", x, a[1]; }' $(MAKEFILE_LIST) | sort

.PHONY: build
## Builds binary from source
build:
	CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "$(LD_FLAGS)" -o dist/dirstalk cmd/dirstalk/main.go
