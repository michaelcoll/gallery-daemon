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
	go run . -f ~/Images/Photos

build:
	go build -v ./...

sqlc:
	sqlc generate
