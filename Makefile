build:
	go build -v -ldflags="-s -w -X 'github.com/michaelcoll/gallery-daemon/cmd.version=v0.0.0'" .

build-docker:
	docker build . -t daemon:latest --build-arg VERSION=latest-local --pull

.PHONY: test
test:
	go test -v ./...

gen: sqlc

clean:
	rm proto/*.pb.go

.PHONY: sqlc
sqlc:
	sqlc generate \
    	&& sqlc-addon generate --quiet

dep-upgrade:
	go get -u
	go mod tidy

docker-run:
	docker run --rm -it -p 9000:9000 daemon:latest
