# Variables
BINARY_NAME=gopher-watch
CONFIG_PATH=configs/targets.json
GO_FILES=./cmd/gopher-watch

.PHONY: all check build run clean errors summary clean-logs stats test

# Default target
all: check build run

# 1. Pre-flight checks (Updated to include test check)
check:
	@echo "Running pre-flight checks..."
	@bash scripts/check_env.sh 
	@echo "All checks passed. Proceeding to build and run."

# 2. Build the binary
build:
	@echo "Building the application..."
	@go build -o ${BINARY_NAME} ${GO_FILES}
	@echo "Build completed successfully."

# 3. Run the application
run: build
	@echo "Running the application..."
	@./${BINARY_NAME}

# 4. NEW: Run Unit Tests
# Includes race detector and coverage report
test:
	@echo "Running unit tests with race detector..."
	@go test -v -race -cover ./internal/...
	@echo "All tests passed."

# 5. Clean up binary
clean:
	@echo "Cleaning up binary..."
	@rm -f ./${BINARY_NAME}
	@echo "Cleanup completed."

# --- Log Management Targets ---

# Filter logs to show only errors
errors:
	@echo "Showing only ERROR logs:"
	@cat gopher-watch.log | jq 'select(.level == "ERROR")'

# Summarize recent probe results
summary:
	@echo "Last 20 probe results:"
	@cat gopher-watch.log | tail -n 20 | jq -r '"[\(.time)] \(.level): \(.target) - \(.msg // "OK")"'

# Clean up the log file
clean-logs:
	@echo "Deleting log file..."
	@rm -f gopher-watch.log

# Provide a quick success/fail count from the logs
stats:
	@echo "Log Statistics Summary:"
	@printf "   Total Probes:  "
	@cat gopher-watch.log | jq '.level' | wc -l
	@printf "   Passed:        "
	@cat gopher-watch.log | jq 'select(.level == "INFO")' | grep -c "level" || echo 0
	@printf "   Failed:        "
	@cat gopher-watch.log | jq 'select(.level == "ERROR")' | grep -c "level" || echo 0