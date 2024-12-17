include .env

.PHONY: list
list:
	@LC_ALL=C $(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/(^|\n)# Files(\n|$$)/,/(^|\n)# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | grep -E -v -e '^[^[:alnum:]]' -e '^$@$$'

stop_containers:
	@echo "Stopping other docker containers..."
	if [ $$(docker ps -q) ]; then \
		echo "Found and stopped containers"; \
		docker stop $$(docker ps -q); \
	else \
		echo "No containers running"; \
	fi

create_container:
	docker run --name ${DB_CONTAINER_NAME} -d \
		-p 3306:3306 \
		-e MYSQL_ROOT_PASSWORD=${DB_ROOT_PASSWORD} \
		-v mysql:/var/lib/mysql \
		mysql:5.7

create_container_db:
	@echo "Creating database ${DB_DATABASE} in Docker container..."
	@if [ -z "$(DB_DATABASE)" ]; then \
		echo "Error: DB_DATABASE is not set. Define it in the .env file."; \
		exit 1; \
	fi
	docker exec -it ${DB_CONTAINER_NAME} mysql -uroot -p${DB_ROOT_PASSWORD} -e "CREATE DATABASE ${DB_DATABASE};"
	@echo "Database ${DB_DATABASE} in Docker container created successfully."

drop_container_db:
	@echo "Dropping database ${DB_DATABASE} in Docker container..."
	@if [ -z "$(DB_DATABASE)" ]; then \
		echo "Error: DB_DATABASE is not set. Define it in the .env file."; \
		exit 1; \
	fi
	docker exec -it ${DB_CONTAINER_NAME} mysql -uroot -p${DB_ROOT_PASSWORD} -e "DROP DATABASE IF EXISTS ${DB_DATABASE};"
	@echo "Database ${DB_DATABASE} in Docker container dropped successfully."

start_container:
	docker start ${DB_CONTAINER_NAME}

create_local_db:
	@echo "Creating local database ${DB_DATABASE}..."
	@if [ -z "$(DB_DATABASE)" ]; then \
		echo "Error: DB_DATABASE is not set. Define it in the .env file."; \
		exit 1; \
	fi
	mysql -u${DB_USERNAME} -p${DB_PASSWORD} -e "CREATE DATABASE ${DB_DATABASE};"
	@echo "Database ${DB_DATABASE} created successfully."

drop_local_db:
	@echo "Dropping local database ${DB_DATABASE}..."
	@if [ -z "$(DB_DATABASE)" ]; then \
		echo "Error: DB_DATABASE is not set. Define it in the .env file."; \
		exit 1; \
	fi
	mysql -u${DB_USERNAME} -p${DB_PASSWORD} -e "DROP DATABASE IF EXISTS ${DB_DATABASE};"
	@echo "Database ${DB_DATABASE} dropped successfully."

create_migrations:
	@if [ -z "$(name)" ]; then \
		echo "Error: You must specify a migration name. Usage: make create_migrations name=your_migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration files for $(name)..."
	migrate create -ext sql -dir ${MIGRATIONS_DIR} $(name)

migrate_up:
	@echo "Running migrations (up)..."
	migrate -path ${MIGRATIONS_DIR} -database "mysql://${DB_USERNAME}:${DB_ROOT_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}" -verbose up

migrate_down:
	@echo "Running migrations (down)..."
	migrate -path ${MIGRATIONS_DIR} -database "mysql://${DB_USERNAME}:${DB_ROOT_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}" -verbose down

build:
	if [ -f "${BINARY}" ]; then \
		rm ${BINARY}; \
		echo "Removed existing binary"; \
	fi
	@echo "Building binary..."
	go build -o ${BINARY} cmd/api/*.go

run:
	go run cmd/api/main.go

build_then_run: build
	@echo "Starting build..."
	./${BINARY}

stop:
	@echo "Stopping server..."
	@pkill -f ${BINARY}
	@echo "Server stopped"