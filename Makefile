build:
	go build -v .

build-prod:
	go build -v -ldflags="-s -w -X 'github.com/michaelcoll/gallery-daemon/cmd.Version=v0.0.0'" .

gen: sqlc protoc

protoc:
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		proto/gallery.proto

clean:
	rm proto/*.pb.go

run:
	go run . index -f ~/Images/Photos

sqlc:
	sqlc generate

dep-upgrade:
	go get -u
	go mod tidy
