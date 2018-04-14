.PHONY: all
all: deps browserpass test

.PHONY: deps
deps:
	dep ensure

browserpass: *.go **/*.go
	go build -o $@

browserpass-linux64: *.go **/*.go
	env GOOS=linux GOARCH=amd64 go build -o $@

browserpass-windows64: *.go **/*.go
	env GOOS=windows GOARCH=amd64 go build -o $@.exe

browserpass-darwinx64: *.go **/*.go
	env GOOS=darwin GOARCH=amd64 go build -o $@

browserpass-openbsd64: *.go **/*.go
	env GOOS=openbsd GOARCH=amd64 go build -o $@

browserpass-freebsd64: *.go **/*.go
	env GOOS=freebsd GOARCH=amd64 go build -o $@

.PHONY: test
test:
	go test
