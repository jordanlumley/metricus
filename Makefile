run-client: 
	@echo "Running client"
	go run cmd/example/client.go

run-agent: 
	@echo "Running agent"
	go run cmd/agent/main.go

start-example-client: 
	@echo "Running example client"
	docker compose up -d --build example_client

start-test-log-container: 
	@echo "Running test log container"
	docker compose up -d --build test_logs

start-agent-container: 
	@echo "Running agent container"
	docker compose up -d --build agent

start-local: 
	@echo "Running local"
	docker compose up -d --build agent test_logs example_client