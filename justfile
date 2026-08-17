# List available recipes
[group('help')]
default:
    @just --list

# Run all unit tests, examples, and saved fuzz regression cases
[group('quality')]
test:
    go test ./...

# Run all tests with Go's race detector
[group('quality')]
test-race:
    go test -race ./...

# Generate and report statement coverage
[group('quality')]
coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out

# Run formatting checks, vet, and tests
[group('quality')]
check: fmt-check vet test

# Format Go code
[group('quality')]
fmt:
    go fmt ./...

# Check that Go code is formatted
[group('quality')]
fmt-check:
    @files="$(gofmt -l $(rg --files -g '*.go'))"; if [ -n "$files" ]; then printf '%s\n' "$files"; exit 1; fi

# Run go vet
[group('quality')]
vet:
    go vet ./...
