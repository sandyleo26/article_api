NAME := blueco-database
VERSION := ${shell git describe --always --long --dirty}

default: container-start

.PHONY: setup
setup:
	go get -u github.com/pressly/goose/cmd/goose
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/gorilla/mux

.PHONY: dep
dep:
	@dep ensure -v

.PHONY: container
container:
	docker build -t test/${NAME}:${VERSION} .

.PHONY: container-start
container-start: container
	docker run -d -p 5432:5432 --name ${NAME} test/${NAME}:${VERSION}

.PHONY: container-stop
container-stop:
	docker rm -f ${NAME}

.PHONY: db-up
db-up:
	goose -dir db/migrations postgres "host=localhost port=5432 dbname=blueco user=postgres password=pass sslmode=disable" up

.PHONY: db-down
db-down:
	goose -dir db/migrations postgres "host=localhost port=5432 dbname=blueco user=postgres password=pass sslmode=disable" down
