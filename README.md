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
      MONGO_URI: mongodb://root:REPLACEME@db:27017  # IMPORTANT: use same credentials as in the <db> container.
      API_URL: http://api:8080
      API_RATE_LIMIT_INTERVAL: 10s
      API_RATE_LIMIT_COUNT: 3
      WORKER_INTERVAL: 5m  # How frequent the API is queried. Should be no less than API update interval below.
    networks:
      - default

  api:
    build:
      # Needs to be built from GitHub as there is currently no public Docker image available
      context: https://github.com/helldivers-2/api.git#ec1e49af15e57fa9b6a464ca5463f3618ee01dac
      dockerfile: ./src/Helldivers-2-API/Dockerfile
    networks:
      - default
    environment:
      Helldivers__Synchronization__IntervalSeconds: 300  # How frequent the API data is updated.
      Helldivers__Synchronization__DefaultLanguage: en-US  # Language of strings such as Major Order text.

  db:
    image: mongo:7.0.7
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: REPLACEME
    networks:
      - default
```
