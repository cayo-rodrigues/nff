DB_NAME := nff-db

TAILWIND_INPUT := ./static/css/custom.css
TAILWIND_OUTPUT := ./static/css/tailwind.css

EC2_HOST := ec2-15-229-157-86.sa-east-1.compute.amazonaws.com

ec2:
	ssh -i "nff-web-key-pair.pem" ubuntu@${EC2_HOST}

turso:
	turso auth login
	turso db shell $(DB_NAME)

tailwind:
	tailwindcss -i $(TAILWIND_INPUT) -o $(TAILWIND_OUTPUT) -w -m
