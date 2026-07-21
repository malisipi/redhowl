live_server:
	find . -name "*" | entr -r go run ./cmd/server