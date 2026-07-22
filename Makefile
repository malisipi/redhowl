live_server:
	find . -name "*" | entr -cr go run ./cmd/server

live_agent:
	find . -name "*" | entr -cr go run ./cmd/agent