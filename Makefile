.PHONY: test local

local:
	docker-compose -f docker/docker-compose.yml up --build
test:
	@echo "Running test"
	ENV=testing go test -v
