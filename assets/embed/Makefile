REGISTRY = node-647ee1368442ecd1a315c673.ps-xaas.io/pluscontainer
IMG ?= {{.ProjectName}}
VERSION ?= $(shell cat VERSION)

.PHONY: embed lint test e2etest build run docker-dev docker-prod docker-push release

#goreleaser release --snapshot --clean

embed:
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		go-bindata -o pkg/project/assets.go -pkg=project -prefix "assets/template"  assets/template/...; \
	else \
		docker-compose run --rm ${IMG} go-bindata -o pkg/assets/embed.go -pkg=assets -prefix "assets/embed"  assets/embed/...; \
	fi

lint: embed
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		go mod tidy && golangci-lint run; \
	else \
		docker-compose run --rm ${IMG} sh -c "go mod tidy && golangci-lint run"; \
	fi

test: embed
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		go mod tidy && go test ./pkg/...; \
	else \
		docker-compose run --rm ${IMG} sh -c "go mod tidy && go test ./pkg/..."; \
	fi

e2etest: embed
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		go mod tidy && go test ./test/...; \
	else \
		docker-compose run --rm ${IMG} sh -c "go mod tidy && go test ./test/..."; \
	fi


build: embed
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		goreleaser build --clean --snapshot; \
	else \
		docker-compose run --rm ${IMG} goreleaser build --clean --snapshot; \
	fi

run:
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		bin/$(IMG); \
	else \
		docker-compose up; \
	fi

docker-dev:
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		echo "Cannot build devcontainer in devcontainer"; \
	else \
		docker-compose build; \
	fi

docker-prod:
	docker build --target prod --platform linux/amd64 -f Dockerfile . -t ${REGISTRY}/${IMG}:${VERSION} -t ${REGISTRY}/${IMG}:latest

docker-push:
	docker push ${REGISTRY}/${IMG}:${VERSION}
	docker push ${REGISTRY}/${IMG}:latest

release:
	docker-compose run ${IMG} goreleaser release --clean