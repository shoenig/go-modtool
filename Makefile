NAME = go-modtool

.PHONY: build
build: clean
	CGO_ENABLED=0 go build -o output/$(NAME)

.PHONY: clean
clean:
	rm -rf dist output/$(NAME)

.PHONY: copywrite
copywrite:
	copywrite headers --spdx "MPL-2.0"

.PHONY: test
test:
	go test -race ./...

.PHONY: test-e2e
test-e2e: clean build
	@echo "[e2e] checking fmt output ..."
	output/$(NAME) -config=e2e/fmt/config.toml fmt e2e/fmt/input.mod > /tmp/fmt.mod
	diff /tmp/fmt.mod e2e/fmt/exp.mod

.PHONY: vet
vet:
	go vet ./...

.PHONY: release
release:
	envy exec gh-release goreleaser release --clean
	$(MAKE) clean

default: build
