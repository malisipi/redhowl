live_server:
	find . -name "*" | entr -r go run ./cmd/server

live_agent:
	find . -name "*" | entr -r go run ./cmd/agent