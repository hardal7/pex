APP_NAME=pex
BUILD_DIR=bin

.PHONY: run-agent run-c2 clean

run-agent:
	go run ./cmd/agent

run-c2:
	go run ./cmd/c2

clean:
	rm -rf $(BUILD_DIR)
