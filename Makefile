SRC_DIRS=cmd pkg

test:
	@echo "Executing tests"
	@CGO_ENABLED=1 go test -v -race -cover ./...

fix:
	@goimports -w $(SRC_DIRS)