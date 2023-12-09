# Define the default target, which is executed when you run just 'make' without any arguments.
.PHONY: default
default: help

# Display help message
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make build         Build Docker images"
	@echo "  make up            Start services in the background"
	@echo "  make down          Stop and remove services"
	@echo ""
	@echo "Variables:"
	@echo "  COMPOSE_FILE      Specify the docker-compose file (default: docker-compose.yml)"
	@echo "  SERVICE           Specify a specific service to build or run tests on"
	@echo "  BUILD_ARGS        Specify additional build arguments"

# Build Docker images
.PHONY: build
build:
	docker-compose -f $(COMPOSE_FILE) build $(SERVICE)

# Start services in the background
.PHONY: up
server:
	docker-compose up -d

# Stop and remove services
.PHONY: down
down:
	docker-compose down
