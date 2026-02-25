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