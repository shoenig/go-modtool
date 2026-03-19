set shell := ["bash", "-u", "-c"]

export scripts := '.github/workflows/scripts'
export GOBIN := `echo $PWD/.bin`

# show available commands
[private]
default:
    @just --list

# compile the executable
[group('build')]
compile: tidy
    go install

# tidy up Go modules
[group('build')]
tidy:
    go mod tidy

# vet the source tree
[group('build')]
vet:
    go vet ./...

# run go test on the source tree
[group('build')]
tests: compile
    go test -race ./...
    echo "[e2e] checking fmt output ..."
    $GOBIN/go-modtool -config=e2e/fmt/config.toml fmt e2e/fmt/input.mod > /tmp/fmt.mod
    diff /tmp/fmt.mod e2e/fmt/exp.mod

# lint the source tree
[group('build')]
lint: vet
    $GOBIN/golangci-lint --config $scripts/golangci.yaml run

# show host system information
[group('setup')]
@sysinfo:
    echo "{{os()/arch()}} {{num_cpus()}}c"

# locally install build dependencies
[group('setup')]
init:
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.1

# create release binaries
[group('release')]
release:
    envy exec gh-release goreleaser release --clean
