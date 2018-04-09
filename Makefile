all: deps browserpass test

.PHONY: deps
deps:
	dep ensure

browserpass: *.go
	go build -o $@

browserpass-linux64: *.go
	env GOOS=linux GOARCH=amd64 go build -o $@

browserpass-windows64: *.go
	env GOOS=windows GOARCH=amd64 go build -o $@.exe

browserpass-darwinx64: *.go
	env GOOS=darwin GOARCH=amd64 go build -o $@

browserpass-openbsd64: *.go
	env GOOS=openbsd GOARCH=amd64 go build -o $@

browserpass-freebsd64: *.go
	env GOOS=freebsd GOARCH=amd64 go build -o $@

.PHONY: test
test:
	go test
