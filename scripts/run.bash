#!/bin/bash

command=$1

if [ -z "$command" ]; then
    command="start"
fi

ProjectRoot="$(dirname "$(readlink -f "$0")")/.."

export AMBULANCE_API_ENVIRONMENT="Development"
export AMBULANCE_API_PORT="8080"
export AMBULANCE_API_MONGODB_USERNAME="root"
export AMBULANCE_API_MONGODB_PASSWORD="neUhaDnes"

function mongo {
    docker compose --file "${ProjectRoot}/deployments/docker-compose/compose.yaml" "$@"
}

case "$command" in
    "openapi")
        docker run --rm -ti -v "${ProjectRoot}:/local" openapitools/openapi-generator-cli generate -c "/local/scripts/generator-cfg.yaml"
        ;;
    "start")
        trap 'mongo down' EXIT
        mongo up --detach
        go run "${ProjectRoot}/cmd/ambulance-api-service"
        ;;
    "test")
        go test -v ./...
        ;;
    "docker")
       docker build -t rbalcercik/ambulance-wl-webapi:local-build -f ${ProjectRoot}/build/docker/Dockerfile .
       ;;
    "mongo")
        mongo up
        ;;
    *)
        echo "Unknown command: $command"
        exit 1
        ;;
esac
