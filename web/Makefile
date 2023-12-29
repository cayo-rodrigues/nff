USER := postgres
PWD :=
ENDPOINT := nff-db.ckwjgza1fp4o.sa-east-1.rds.amazonaws.com
DB_NAME := nff

DIR_PATH := internal/db/migrations
NAME :=

db_prod:
	docker run -it --rm postgres psql "postgres://$(USER):$(PWD)@$(ENDPOINT):5432/$(DB_NAME)"

new_migration:
	dc run migrations -e CMD='create $(NAME) sql'
	chmod o+w $(DIR_PATH)/$(NAME).sql
