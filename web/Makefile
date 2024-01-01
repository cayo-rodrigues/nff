USER := postgres
PWD :=
ENDPOINT := nff-db.ckwjgza1fp4o.sa-east-1.rds.amazonaws.com
DB_NAME := nff

DIR_PATH := db/migrations
NAME :=

rds:
	docker run -it --rm -e PGPASSWORD=$(PWD) postgres psql "postgres://$(USER)@$(ENDPOINT):5432/$(DB_NAME)"

new_migration:
	docker compose run -e CMD='create $(NAME) sql' migrations
	sudo chmod o+w $(call get_latest_migration_file)

define get_latest_migration_file
	$(DIR_PATH)/$$(ls -lt $(DIR_PATH) | awk 'NR==2{print $$9}')
endef
