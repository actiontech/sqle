GIT_VERSION   = $(shell git rev-parse --abbrev-ref HEAD) $(shell git rev-parse HEAD)
LDFLAGS       = -ldflags "-X 'main.version=\"${GIT_VERSION}\"'"
RPM_BUILD_BIN = $(shell type -p rpmbuild 2>/dev/null)
COMPILE_FLAG  =
DOCKER        = $(shell which docker)
DOCKER_IMAGE  = docker-registry:5000/actiontech/universe-compiler-go1.10
DOCKER_REGISTRY ?=10.186.18.20

PROJECT_NAME = sqle
VERSION       = 9.9.9.9

default: build

pull_image:
    $(DOCKER) pull ${DOCKER_IMAGE}

build:
	go build -o sqled -ldflags "-X 'main.version=\"${GIT_VERSION}\"'" ./cmd/main.go

vet:
	go vet $$(GOPATH=${GOPATH} GOOS=linux GOARCH=amd64 go list ./... | grep -v vendor/)

test:
	go test -v ./...

clean:
	GOPATH=${GOPATH} GOOS=linux GOARCH=amd64 go clean

docker_rpm: pull_image
	$(DOCKER) run -v $(shell pwd):/universe/src/sqle --rm $(DOCKER_IMAGE) -c "(mkdir -p /root/rpmbuild/SOURCES >/dev/null 2>&1); cd /root/rpmbuild/SOURCES; (tar zcf ${PROJECT_NAME}.tar.gz /universe --exclude=/universe/pkg --exclude=/universe/bin --transform 's/universe/${PROJECT_NAME}-${VERSION}/' >/tmp/build.log 2>&1) && (rpmbuild -bb --with qa /universe/src/sqle/build/sqled.spec >>/tmp/build.log 2>&1) && (cat /root/rpmbuild/RPMS/x86_64/${PROJECT_NAME}-${VERSION}-qa.x86_64.rpm) || (cat /tmp/build.log && exit 1)" > ${PROJECT_NAME}.x86_64.rpm

docker_test: pull_image
	CTN_NAME="universe_docker_test_$$RANDOM" && \
    $(DOCKER) run -d --entrypoint /sbin/init --add-host docker-registry:${DOCKER_REGISTRY}  --privileged --name $${CTN_NAME} -v $(shell pwd)/../:/universe --rm -w /universe/src $(DOCKER_IMAGE) && \
    $(DOCKER) exec $${CTN_NAME} make test ; \
    $(DOCKER) stop $${CTN_NAME}

upload:
	curl -T $(shell pwd)/${PROJECT_NAME}.x86_64.rpm -u admin:ftpadmin ftp://release-ftpd/actiontech-${PROJECT_NAME}/qa/${VERSION}/${PROJECT_NAME}-${VERSION}-qa.x86_64.rpm
