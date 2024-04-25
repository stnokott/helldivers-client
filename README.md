# Helldivers 2 Client

> [!IMPORTANT]  
> Currently in very active development, please don't use unless you are ready for chaos!

## Summary

### Goal
 Storing metrics as historic data for later processing (i.e. visualization)

### Procedure
1. Query the Helldivers 2 API at regular intervals
2. Write the received data into snapshots, each associated with the current time

## Setup

Use the following `docker-compose.yaml` to get started:

```yaml
version: '3.8'

networks:
  default:
    driver: bridge

services:
  app:
    image: ghcr.io/stnokott/helldivers-client:latest
    depends_on:
      - api
      - db
    environment:
      MONGO_URI: postgresql://root:REPLACEME@db:5432/helldivers  # IMPORTANT: use same credentials as in the <db> container.
      API_URL: http://api:8080
      WORKER_INTERVAL: 5m  # How frequent the API is queried. Should be no less than API update interval below.
    networks:
      - default

  api:
    build:
      # Needs to be built from GitHub as there is currently no public Docker image available
      context: https://github.com/helldivers-2/api.git#f754f2bb8d53e278e4239ff08d6daa8cb803db66
      dockerfile: ./src/Helldivers-2-API/Dockerfile
    networks:
      - default
    environment:
      Helldivers__Synchronization__IntervalSeconds: 300  # How frequent the API data is updated.
      Helldivers__Synchronization__DefaultLanguage: en-US  # Language of strings such as Major Order text.

  db:
    image: postgres:16.2
    environment:
      POSTGRES_USER: root
      POSTGRES_PASSWORD: REPLACEME
      POSTGRES_DB: helldivers
    volumes:
      - /path/to/persistent/storage:/var/lib/postgresql/data  # Persist your DB data
    networks:
      - default
```

## Development

### Install required tools

```shell
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0
go install github.com/joho/godotenv/cmd/godotenv@v1.5.1
go install github.com/goreleaser/goreleaser@v1.25.1
go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0
```
