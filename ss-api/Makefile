test:
	docker run --rm -it -e CHROMEDRIVER_PATH=/opt/chromedriver cayy/nff-ss-api-dev python tests.py $(HANDLER)

build_dev:
	docker build -t cayy/nff-ss-api-dev -f Dockerfile.dev .
