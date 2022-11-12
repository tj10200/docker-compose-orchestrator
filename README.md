# Docker Compose Orchestrator

This tool utilizes existing docker-compose `depends_on` keywords to build
an execution graph. Additional configuration is then supplied to run commands
to verify readiness for the service with a given dependency.

To run:

    go run main.go \
    --docker-compose-files "/path/to/docker-compose.yml" \
    --config "./docker-compose-orchestrator/config.yaml" \
    --env "/path/to/.env"

See the [Example configuration](example/example_config.yaml) to get `config.yaml` structure