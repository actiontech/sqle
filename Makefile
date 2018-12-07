GOCMD=$(shell which go)
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOLIST=$(GOCMD) list
GOCLEAN=$(GOCMD) clean

GIT_VERSION   = $(shell git rev-parse --abbrev-ref HEAD) $(shell git rev-parse HEAD)
LDFLAGS       = -ldflags "-X 'main.version=\"${GIT_VERSION}\"'"
RPM_BUILD_BIN = $(shell type -p rpmbuild 2>/dev/null)
COMPILE_FLAG  =
DOCKER        = $(shell which docker)
DOCKER_IMAGE  = docker-registry:5000/actiontech/universe-compiler-go1.10
DOTNET_DOCKER_IMAGE = docker-registry:5000/actiontech/universe-compiler-dotnetcore2.1
DOCKER_REGISTRY ?=10.186.18.20
DOTNET_TARGET = centos.7-x64

PROJECT_NAME = sqle
SUB_PROJECT_NAME = sqle_sqlserver
VERSION       = 9.9.9.9

.PHONY: build docs

default: build

pull_image:
    $(DOCKER) pull ${DOCKER_IMAGE}

build: swagger parser vet
	$(GOBUILD) -o sqled -ldflags "-X 'main.version=\"${GIT_VERSION}\"'" ./cmd/main.go

build_sqlserver:
	cd ./sqlserver/SqlserverProtoServer && dotnet publish -c Release -r ${DOTNET_TARGET}

vet: swagger
	$(GOVET) $$($(GOLIST) ./... | grep -v vendor/)

test: swagger parser
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)

docker_rpm: pull_image
	$(DOCKER) run -v $(shell pwd):/universe/src/sqle --rm $(DOCKER_IMAGE) -c "(mkdir -p /root/rpmbuild/SOURCES >/dev/null 2>&1);cd /root/rpmbuild/SOURCES; (tar zcf ${PROJECT_NAME}.tar.gz /universe --transform 's/universe/${PROJECT_NAME}-${VERSION}/' >/tmp/build.log 2>&1) && (rpmbuild -bb --with qa /universe/src/sqle/build/sqled.spec >>/tmp/build.log 2>&1) && (cat /root/rpmbuild/RPMS/x86_64/${PROJECT_NAME}-${VERSION}-qa.x86_64.rpm) || (cat /tmp/build.log && exit 1)" > ${PROJECT_NAME}.x86_64.rpm
	$(DOCKER) run -v $(shell pwd):/universe/src/sqle --rm $(DOTNET_DOCKER_IMAGE) -c "(mkdir -p /root/rpmbuild/SOURCES >/dev/null 2>&1);cd /root/rpmbuild/SOURCES; (tar zcf ${SUB_PROJECT_NAME}.tar.gz /universe --transform 's/universe/${SUB_PROJECT_NAME}-${VERSION}/' >/tmp/build.log 2>&1) && (DOTNET_TARGET=${DOTNET_TARGET} rpmbuild -bb --with qa /universe/src/sqle/build/sqled_sqlserver.spec >>/tmp/build.log 2>&1) && (cat /root/rpmbuild/RPMS/x86_64/${SUB_PROJECT_NAME}-${VERSION}-qa.x86_64.rpm) || (cat /tmp/build.log && exit 1)" > ${SUB_PROJECT_NAME}.x86_64.rpm

docker_test: pull_image
	CTN_NAME="universe_docker_test_$$RANDOM" && \
    $(DOCKER) run -d --entrypoint /sbin/init --add-host docker-registry:${DOCKER_REGISTRY}  --privileged --name $${CTN_NAME} -v $(shell pwd):/universe/src/sqle --rm -w /universe/src/sqle $(DOCKER_IMAGE) && \
    $(DOCKER) exec $${CTN_NAME} make test ; \
    $(DOCKER) stop $${CTN_NAME}

upload:
	curl -T $(shell pwd)/${PROJECT_NAME}.x86_64.rpm -u admin:ftpadmin ftp://release-ftpd/actiontech-${PROJECT_NAME}/qa/${VERSION}/${PROJECT_NAME}-${VERSION}-qa.x86_64.rpm

parser:
	cd $(shell pwd)/vendor/github.com/pingcap/tidb && make parser && cd -

swagger:
	$(GOBUILD) -o $(shell pwd)/swag $(shell pwd)/vendor/github.com/swaggo/swag/cmd/swag/main.go
	$(shell pwd)/swag init -g $(shell pwd)/api/app.go