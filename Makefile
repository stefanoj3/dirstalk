SRC_DIRS=cmd pkg

TESTARGS=-v -race -cover

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

.PHONY: fix
## Run goimports against the source code
fix:
	@goimports -w $(SRC_DIRS)

.PHONY: release-snapshot
## Creates a release snapshot - requires goreleaser to be available in the $PATH
release-snapshot:
	@echo "Creating release snapshot..."
	@goreleaser --snapshot --rm-dist

.PHONY: help
## Display this help screen - requires gawk
help:
	@gawk 'match($$0, /^## (.*)/, a) \
		{ getline x; x = gensub(/(.+:) .+/, "\\1", "g", x) ; \
		printf "\033[36m%-30s\033[0m %s\n", x, a[1]; }' $(MAKEFILE_LIST) | sort
