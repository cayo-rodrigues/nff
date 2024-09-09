TURSO_DATABASE_URL := 
TURSO_AUTH_TOKEN := 

MIGRATIONS_DIR := sql
NAME :=

build:
	go build -o ./bin/migrate

run: build
	@TURSO_DATABASE_URL=$(TURSO_DATABASE_URL) \
	TURSO_AUTH_TOKEN=$(TURSO_AUTH_TOKEN) \
	./bin/migrate $(call get_cmd_args)

new_migration:
	$(MAKE) run create $(NAME)
	sudo chmod o+w $(call get_latest_migration_file)

define get_latest_migration_file
	$(MIGRATIONS_DIR)/$$(ls -lt $(MIGRATIONS_DIR) | awk 'NR==2{print $$9}')
endef

define get_cmd_args
	$(filter-out $@,$(MAKECMDGOALS))
endef

%:
	@:
