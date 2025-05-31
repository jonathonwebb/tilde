BINDIR := bin
TMPDIR := tmp

.PHONY: build
build:
	CGO_ENABLED=1 go build -o ./bin/tilde ./cmd

.PHONY: check
check: deps-check fmt-check lint

.PHONY: deps
deps:
	go mod tidy -v

.PHONY: deps-check
deps-check:
	go mod tidy -diff
	go mod verify

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: fmt-check
fmt-check:
	test -z "$(shell gofmt -l .)"

.PHONY: lint
lint:
	golangci-lint run

.PHONY: vet
	go vet ./...

.PHONY: test
test:
	go test ./...

.PHONY: test-check
test-check:
	go test -race -count=1 ./...

.PHONY: cover
cover: $(TMPDIR)
	go test -v -coverprofile $(TMPDIR)/cover.out ./...
	go tool cover -html=$(TMPDIR)/cover.out

.PHONY: cover-check
cover-check: $(TMPDIR)
	go test -race -count=1 -coverprofile $(TMPDIR)/cover.out ./...

.PHONY: clean
clean:
	rm -rf $(BINDIR) $(TMPDIR)

$(TMPDIR) $(BINDIR):
	mkdir -p $@
