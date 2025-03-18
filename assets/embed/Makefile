REGISTRY = node-647ee1368442ecd1a315c673.ps-xaas.io/pluscontainer
OS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH ?= $(shell uname -m)
ifeq ($(ARCH),aarch64)
  ARCH := arm64
endif
VERSION ?= $(shell cat VERSION)
IMG ?= {{.ProjectName}}

.PHONY: embed lint test e2etest build run release docker-dev docker-prod docker-release

embed:
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		go-bindata -o pkg/project/assets.go -pkg=project -prefix "assets/template"  assets/template/...; \
	else \
		docker-compose run -T --rm --remove-orphans ${IMG} go-bindata -o pkg/assets/embed.go -pkg=assets -prefix "assets/embed"  assets/embed/...; \
	fi

lint: embed
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		go mod tidy && golangci-lint run; \
	else \
		docker-compose run -T --rm --remove-orphans ${IMG} sh -c "go mod tidy && golangci-lint run"; \
	fi

test: lint
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		go test ./pkg/...; \
	else \
		docker-compose run -T --rm --remove-orphans ${IMG} go test ./pkg/...; \
	fi

e2etest: lint
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		go test ./test/...; \
	else \
		docker-compose run -T --rm --remove-orphans ${IMG} go test ./test/...; \
	fi


build: embed
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		goreleaser build --clean --snapshot; \
	else \
		docker-compose run -T --rm --remove-orphans ${IMG} goreleaser build --clean --snapshot; \
	fi

run: embed
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		dist/$(IMG)_${OS}_${ARCH}*/${IMG}; \
	else \
		docker-compose up; \
	fi

release: embed
	docker-compose run -T --rm --remove-orphans ${IMG} goreleaser release --clean

docker-dev: 
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		echo "Cannot build devcontainer in devcontainer"; \
	else \
		docker-compose build; \
	fi

docker-prod: 
	docker build --target prod -f Dockerfile . -t ${REGISTRY}/${IMG}:${VERSION} -t ${REGISTRY}/${IMG}:latest

docker-release:
	docker build --target prod --platform linux/amd64 -f Dockerfile . -t ${REGISTRY}/${IMG}:${VERSION} -t ${REGISTRY}/${IMG}:latest
	docker push --platform linux/amd64 ${REGISTRY}/${IMG}:${VERSION}
	docker push --platform linux/amd64 ${REGISTRY}/${IMG}:latest
