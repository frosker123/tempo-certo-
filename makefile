.PHONY: build
build:
	@go build ./cmd/api

.PHONY: run
run:
	@go run ./cmd/api	

.PHONY: docker-stop
docker-stop:
	@docker-compose stop

.PHONY: docker
docker: 
	@docker-compose up 

