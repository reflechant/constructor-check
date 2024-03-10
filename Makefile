constructor-check: main.go analyzer/analyzer.go
	go build -o constructor-check

.PHONY: build
build: constructor-check

.PHONY: test
test:
	go test ./analyzer