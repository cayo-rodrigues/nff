USER := postgres
PWD :=
ENDPOINT := nff-db.ckwjgza1fp4o.sa-east-1.rds.amazonaws.com
PORT := 5432
DB_NAME := nff
RDS_URL := postgres://$(USER)@$(ENDPOINT):$(PORT)/$(DB_NAME)
RDS_URL_WITH_PWD := postgres://$(USER):$(PWD)@$(ENDPOINT):$(PORT)/$(DB_NAME)

DIR_PATH := database/migrations
NAME :=

TAILWIND_INPUT := ./static/css/custom.css
TAILWIND_OUTPUT := ./static/css/tailwind.css

EC2_HOST := ec2-52-67-84-134.sa-east-1.compute.amazonaws.com

ec2:
	ssh -i "nff-web-key-pair.pem" ubuntu@${EC2_HOST}

rds:
	docker run -it --rm -e PGPASSWORD=$(PWD) postgres psql $(RDS_URL)

rds_rollback:
	docker compose run -e CMD='down' -e DB_URL=$(RDS_URL_WITH_PWD) migrations

new_migration:
	docker compose run -e CMD='create $(NAME) sql' migrations
	sudo chmod o+w $(call get_latest_migration_file)

tailwind:
	tailwindcss -i $(TAILWIND_INPUT) -o $(TAILWIND_OUTPUT) -w -m

define get_latest_migration_file
	$(DIR_PATH)/$$(ls -lt $(DIR_PATH) | awk 'NR==2{print $$9}')
endef
