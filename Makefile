constructor-check: constructor-check.go cmd/constructor-check/main.go
	go build -o constructor-check cmd/constructor-check/main.go

.PHONY: build
build: constructor-check

.PHONY: test
test:
	go test