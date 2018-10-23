.PHONY: test lint
TEMPDIR := $(shell mktemp -d)

all: install

test: gotest lint

gotest:
	@go test -v --cover ./...

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

lint: $(GOMETALINTER)
	@CGO_ENABLED=0 gometalinter ./...
	gometalinter ./... --vendor

install:
	@CGO_ENABLED=0 go build

serve: docker_build docker_run

docker_build:
	@docker build -t go-simple-server .

docker_run:
	@echo Mounting Tempdir into container: $(TEMPDIR)
	@docker run -ti -p 8081:8081 -v $(TEMPDIR):/workdir/serve go-simple-server

clean:
	@rm -rf ./tmp
	@docker rmi go-simple-server

