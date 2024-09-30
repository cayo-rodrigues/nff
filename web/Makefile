USER := postgres
PWD :=
ENDPOINT := nff-db.ckwjgza1fp4o.sa-east-1.rds.amazonaws.com
PORT := 5432
DB_NAME := nff
RDS_URL := postgres://$(USER)@$(ENDPOINT):$(PORT)/$(DB_NAME)
RDS_URL_WITH_PWD := postgres://$(USER):$(PWD)@$(ENDPOINT):$(PORT)/$(DB_NAME)

TAILWIND_INPUT := ./static/css/custom.css
TAILWIND_OUTPUT := ./static/css/tailwind.css

EC2_HOST := ec2-15-229-157-86.sa-east-1.compute.amazonaws.com

ec2:
	ssh -i "nff-web-key-pair.pem" ubuntu@${EC2_HOST}

turso:
	docker run -it --rm -e PGPASSWORD=$(PWD) postgres psql $(RDS_URL)

tailwind:
	tailwindcss -i $(TAILWIND_INPUT) -o $(TAILWIND_OUTPUT) -w -m
