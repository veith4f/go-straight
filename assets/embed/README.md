# $project_name
This project uses docker-compose (https://docs.docker.com/compose/) to set up a local \
development environment where all dependencies are installed into containers. \
Docker-compose allows defining multiple containers - such as a backend, a database, a frontend \
that can communicate with each other. 

## Getting started
- \`make lint\`: lint all files
- \`make build\`: execute go build
- \`make run\`: execute binary
- \`make compose-build\`: build docker containers
- \`make up\`: run compose containers
- \`make compose-test\`: execute tests in docker container
