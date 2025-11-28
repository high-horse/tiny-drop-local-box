goose-create:
	goose -dir migrations create init sql


goose-up:
	goose -dir migrations sqlite3 storage/files.db up

generate:
	sqlc-generate


build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/tinydrop ./cmd


