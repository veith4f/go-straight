REGISTRY = node-647ee1368442ecd1a315c673.ps-xaas.io/pluscontainer
IMG ?= go-straight
VERSION ?= $(shell cat VERSION)

.PHONY: lint test build run

lint:
	@go mod tidy
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run; \
	else \
		bin/golangci-lint run; \
	fi

test:
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		command -v go &> /dev/null || ( echo "Please install go" && exit 1 ); \
		go mod tidy && go test ./pkg/...; \
	else \
		command -v docker-compose &> /dev/null || ( echo "Please install docker-compose" && exit 1 ); \
		docker-compose run --rm --remove-orphans ${IMG} sh -c "go mod tidy && go test ./pkg/..."; \
	fi

e2etest:
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		command -v go &> /dev/null || ( echo "Please install go" && exit 1 ); \
		go mod tidy && go test ./test/...; \
	else \
		command -v go &> /dev/null || ( echo "Please install go" && exit 1 ); \
		go mod tidy && go test ./test/...; \
	fi


build:
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		command -v go &> /dev/null || ( echo "Please install go" && exit 1 ); \
		go build -o bin/${IMG} cmd/main.go; \
	else \
		command -v go &> /dev/null || ( echo "Please install docker-compose" && exit 1 ); \
		docker-compose build; \
	fi

run:
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		bin/$(IMG); \
	else \
		command -v go &> /dev/null || ( echo "Please install docker-compose" && exit 1 ); \
		docker-compose up; \
	fi

docker-build:
	docker build -f docker/${IMG} . --target prod -t ${REGISTRY}/${IMG}:${VERSION} -t ${REGISTRY}/${IMG}:latest

docker-push:
	docker push ${REGISTRY}/${IMG}:${VERSION}
	docker push ${REGISTRY}/${IMG}:latest
