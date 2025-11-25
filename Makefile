goose-create:
	goose -dir migrations create init sql


goose-up:
	goose -dir migrations sqlite3 storage/files.db up

generate:
	sqlc-generate

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o myapp .
