NAME := blueco-database
VERSION := ${shell git describe --always --long --dirty}

default: start

.PHONY: setup
setup:
	go get -u github.com/pressly/goose/cmd/goose
	go get -u github.com/golang/dep/cmd/dep
	$(MAKE) dep

.PHONY: dep
dep:
	@dep ensure -v

.PHONY: start
start:
	go run main.go

.PHONY: db-container
db-container:
	docker build -t test/${NAME}:${VERSION} .

.PHONY: db-container-start
db-container-start: db-container
	docker run -d -p 5432:5432 --name ${NAME} test/${NAME}:${VERSION}

.PHONY: db-up
db-up:
	goose -dir db/migrations postgres "host=localhost port=5432 dbname=blueco user=postgres password=pass sslmode=disable" up

.PHONY: db-down
db-down:
	goose -dir db/migrations postgres "host=localhost port=5432 dbname=blueco user=postgres password=pass sslmode=disable" down

.PHONY: db-stop
db-stop:
	docker rm -f ${NAME}