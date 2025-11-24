goose-create:
	goose -dir migrations create init sql


goose-up:
	goose -dir migrations sqlite3 storage/files.db up

generate:
	sqlc-generate