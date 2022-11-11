build:
	go build -v -ldflags="-s -w -X 'github.com/michaelcoll/gallery-daemon/cmd.version=v0.0.0'" .

build-docker:
	docker build . -t daemon:latest --build-arg VERSION=latest-local

.PHONY: test
test:
	go test -v ./...

gen: sqlc

clean:
	rm proto/*.pb.go

.PHONY: sqlc
sqlc:
	sqlc generate

dep-upgrade:
	go get -u
	go mod tidy
