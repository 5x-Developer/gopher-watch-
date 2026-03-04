#Variables
BINARY_NAME=gopher-watch
CONFIG_PATH=configs/targets.json

#Commands
all: check run

#Pre Filght check
check:
	@echo "Running pre-flight checks..."
	@bash scripts/check_env.sh 
	@echo "All checks passed. Proceeding to build and run."

#Build and run

build:
	@echo "Building the application..."
	@go build -o ${BINARY_NAME} ./cmd/gopher-watch
	@echo "Build completed successfully."

run: build
	@echo "Running the application..."
	@./${BINARY_NAME}

clean:
	@echo "Cleaning up..."
	@rm -f ./${BINARY_NAME}
	@echo "Cleanup completed."

# Filter logs to show only errors
errors:
	@cat gopher-watch.log | jq 'select(.level == "ERROR")'

# Summarize recent probe results
summary:
	@cat gopher-watch.log | tail -n 20 | jq -r '"[\(.time)] \(.level): \(.target) - \(.msg // "OK")"'

# Clean up the log file
clean-logs:
	rm -f gopher-watch.log

	# Provide a quick success/fail count from the logs
stats:
	@echo "Log Statistics Summary:"
	@printf "   Total Probes:  "
	@cat gopher-watch.log | jq '.level' | wc -l
	@printf "   Passed:      "
	@cat gopher-watch.log | jq 'select(.level == "INFO")' | grep -c "level" || echo 0
	@printf "   Failed:      "
	@cat gopher-watch.log | jq 'select(.level == "ERROR")' | grep -c "level" || echo 0