# launchd_docker

Helper utility to manage docker containers on macOS using Lima VM.

## Features

- Manages Docker services using Docker Compose
- Automatically manages Lima VM lifecycle
- Supports multiple services with independent configurations
- Graceful shutdown handling
- Structured logging
- Custom compose file support

## Prerequisites

- Go 1.21 or later
- Lima VM installed and configured
- Docker and Docker Compose installed
- macOS

## Building

```bash
go build -o launchd_docker cmd/launchd_docker/main.go
```

## Configuration

Create a YAML configuration file with the following structure:

```yaml
hypervisor:
  lima_instance: "default"  # Name of your Lima instance

services:
  - name: "my-app"
    path: "/Users/username/projects/my-app"
    compose_file: "docker-compose.yaml"  # Optional, defaults to docker-compose.yaml
```

## Usage

```bash
./launchd_docker -config path/to/config.yaml
```

The program will:
1. Validate the configuration file
2. Check if the Lima VM is running and start it if needed
3. Start all configured services using Docker Compose
4. Handle graceful shutdown on SIGINT or SIGTERM

## Error Handling

- All errors are logged to stderr
- Service failures are logged but don't stop other services
- Graceful shutdown ensures all services are properly stopped

## License

MIT License 