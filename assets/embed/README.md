# {{.ProjectName}}
This project uses docker-compose (https://docs.docker.com/compose/) to set up a local \
development environment where all dependencies are installed into containers. \
Docker-compose allows defining multiple containers - such as a backend, a database, a frontend \
that can communicate with each other. 

## Working in a project
```makefile
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
```
