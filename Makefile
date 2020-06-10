SNAPSHOT_NUMBER :=$(shell date +'%y%m%d')
VERSION_NUMBER=0.1
VERSION=$(VERSION_NUMBER)-b$(SNAPSHOT_NUMBER)

BINARY = bin
SRC = cmd/amqp-cli
APP = amqp-cli

.DEFAULT_GOAL := all

.PHONY: all
all: clean build

.PHONY: pre-compile
pre-compile:
	@mkdir -p ${BINARY}

compile: pre-compile
	@echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=amd64 go build -o ${BINARY}/${APP}-v${VERSION}-linux-amd64 ./${SRC}
	GOOS=windows GOARCH=amd64 go build -o ${BINARY}/${APP}-v${VERSION}-windows-amd64 ./${SRC}

build: compile

clean:
	@if [ -d ${BINARY} ] ; then rm -rf ${BINARY}/* ; fi

.PHONY: sync
sync:
	scp ${BINARY}/*-linux-amd64 localvm:/home/mostafa/projects/go-amqp-cli
	scp cmd/amqp-cli/*.yaml localvm:/home/mostafa/projects/go-amqp-cli
	ssh localvm ln -s /home/mostafa/projects/go-amqp-cli/${APP}-v${VERSION}-linux-amd64 /home/mostafa/projects/go-amqp-cli/go-amqp-cli
