AUTH_BINARY=auth-service

go_build:
	@echo "Building ${AUTH_BINARY} binary..."
	@GOOS=linux GOARCH=amd64 go build -o ./bin/app/${AUTH_BINARY} ./cmd/app
	@echo "Done!"

clean_build:
	@echo "Removing binaries..."
	@rm -f ./bin/app/${AUTH_BINARY}
	@echo "Done!"

go_run: go_build
	@echo "Running ${AUTH_BINARY} binary..."
	@export CONFIG_PATH=$(CONFIG_PATH) && ./bin/app/${AUTH_BINARY} &
	@echo "Done!"

go_stop:
	@echo "Stopping ${AUTH_BINARY}..."
	@-pkill -SIGTERM -f "./bin/app/${AUTH_BINARY}"
	@echo "Stopped ${AUTH_BINARY}!"

build_image:
	@echo "Building image..."
	@docker build --no-cache -f build/docker/Dockerfile -t ${AUTH_BINARY} .
	@echo "Done!"

run:
	@echo "Running container..."
	@docker run --rm -d --name=${AUTH_BINARY} \
	-v ./config:/usr/local/bin/config \
	-e CONFIG_PATH=/usr/local/bin/config/config.yaml \
	-p6000:6000 ${AUTH_BINARY}
	@echo "Done!"

stop:
	@echo "Stopping container..."
	@docker stop ${AUTH_BINARY}
	@echo "Done!"
