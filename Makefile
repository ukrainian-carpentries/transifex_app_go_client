.PHONY: main go

main: go

go:
	git pull
	go mod tidy
	go build ./...
