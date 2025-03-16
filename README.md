# Go straight
CLI-tool that scaffolds a go project targeted at shipping binaries and or docker containers. Go straight sets you up with:
- code linting using [golangci-lint](https://golangci-lint.run/)
- unit tests using [go testing](https://pkg.go.dev/testing)
- end-to-end tests using using [Ginkgo](https://onsi.github.io/ginkgo/)
- [VS Code Devcontainer](https://code.visualstudio.com/docs/devcontainers/containers)
- containerized builds and shipping using [docker-compose](https://docs.docker.com/compose/)
- baked-in assets using [go-bindata](https://github.com/go-bindata/go-bindata)
- binary releases for all platforms using [goreleaser](https://goreleaser.com/)
- [Makefile](https://www.gnu.org/software/make/) with streamlined operations

## Getting started
```makefile
embed:
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		go-bindata -o pkg/project/assets.go -pkg=project -prefix "assets/template"  assets/template/...; \
	else \
		docker-compose run --rm --remove-orphans ${IMG} go-bindata -o pkg/assets/embed.go -pkg=assets -prefix "assets/embed"  assets/embed/...; \
	fi

lint: embed
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		go mod tidy && golangci-lint run; \
	else \
		docker-compose run --rm --remove-orphans ${IMG} golangci-lint run; \
	fi

test: embed
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		go mod tidy && go test ./pkg/...; \
	else \
		docker-compose run --rm --remove-orphans ${IMG} go test ./pkg/...; \
	fi

e2etest: embed
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		go mod tidy && go test ./test/...; \
	else \
		docker-compose run --rm --remove-orphans ${IMG} go test ./test/...; \
	fi


build: embed
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		goreleaser build --clean --snapshot; \
	else \
		docker-compose run --rm --remove-orphans ${IMG} goreleaser build --clean --snapshot; \
	fi

run: embed
	@if [ -f /.dockerenv ] || ( [ -f /proc/self/cgroup ] && grep -qE 'docker|containerd' /proc/self/cgroup ); then \
		dist/$(IMG)_${OS}_${ARCH}*/${IMG}; \
	else \
		docker-compose up; \
	fi

release: embed
	docker-compose run --rm --remove-orphans ${IMG} goreleaser release --clean

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

```
