SRC_DIRS=cmd pkg

ifeq ($(TRAVIS), true)
TESTARGS=-v -race -coverprofile=coverage.txt -covermode=atomic
else
TESTARGS=-v -race -cover
endif

.PHONY: dep
dep:
	@go get github.com/golang/dep/cmd/dep
	@go get golang.org/x/tools/cmd/goimports
	@go get golang.org/x/lint/golint
	@dep ensure

.PHONY: test
test:
	@echo "Executing tests"
	@CGO_ENABLED=1 go test $(TESTARGS) ./...


.PHONY: check
check:
	@golint -set_exit_status .
	@go vet ./...
	@goimports -l $(SRC_DIRS) | tee /dev/tty | xargs -I {} test -z {}

.PHONY: fix
fix:
	@goimports -w $(SRC_DIRS)
