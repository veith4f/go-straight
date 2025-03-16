FROM docker.io/golang:1.24.1 AS dev

ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV ORIGPATH=$PATH
ENV PATH=$GOBIN:$PATH

RUN echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | tee /etc/apt/sources.list.d/goreleaser.list && \
apt update && apt install -y goreleaser ca-certificates git vim nano curl wget python3-pip pipx unzip docker && \
curl -Lo ./kind https://kind.sigs.k8s.io/dl/latest/kind-linux-amd64 && \
chmod +x ./kind &&mv ./kind /usr/local/bin/kind && \
curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/linux/amd64 && \
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/ && \
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
chmod +x kubectl && mv kubectl /usr/local/bin/kubectl && \
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && \
unzip awscliv2.zip && ./aws/install && rm awscliv2.zip && \
curl -SL https://github.com/docker/compose/releases/download/v2.33.1/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose && \
pipx install python-openstackclient && pipx ensurepath && \
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.7 && \
go install github.com/go-bindata/go-bindata/...@latest

WORKDIR /workspace
COPY . .

RUN go-bindata -o pkg/assets/embed.go -pkg=assets -prefix "assets/embed"  assets/embed/...

ENV GOPATH=/workspace/.go
ENV GOBIN=/workspace/.go/bin
ENV PATH=$GOBIN:$ORIGPATH

RUN go mod tidy && golangci-lint run
RUN go build -o /go-straight cmd/main.go

CMD ["/go-straight"]

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot AS prod
WORKDIR /
COPY --from=dev /go-straight .
USER 65532:65532
