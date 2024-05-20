![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/stnokott/helldivers-client)
![GitHub Release](https://img.shields.io/github/v/release/stnokott/helldivers-client?logo=docker)
![Github Container Registry Image Size](https://ghcr-badge.egpl.dev/stnokott/helldivers-client/size?tag=latest)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/stnokott/helldivers-client/test.yml?branch=main&event=schedule&label=integration%20tests)

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
      POSTGRES_URI: postgresql://root:REPLACEME@db:5432/helldivers  # IMPORTANT: use same credentials as in the <db> container.
      API_URL: http://api:8080
      WORKER_INTERVAL: 5m  # How frequent the API is queried. Should be no less than API update interval below.
      TZ: Europe/Berlin
    networks:
      - default

  api:
    build:
      # Needs to be built from GitHub as there is currently no public Docker image available
      context: https://github.com/helldivers-2/api.git#9c6869c0fa93d41f715fd55b8cce1f11ca257475  # pin version
      dockerfile: ./src/Helldivers-2-API/Dockerfile
    networks:
      - default
    environment:
      Helldivers__API__Authentication__Enabled: false  # Set to true if exposed
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

### PGO

We use `PGO` for build optimization.
To collect a new CPU profile from a representative productive environment, do the following:
1. Run the "Release snapshot" action on your current working branch.
   This will build and publish a snapshot release of the code on your branch **without PGO**.
2. Pull the `ghcr.io/stnokott/helldivers-client:latest-snapshot` Docker image
3. Run this image in the same Docker stack as you would do with the regular image.
   Use the following parameters to control profiling behaviour:

   `--pprof-duration=30m` -> How long the profiling should run for.
   
   `--pprof-out=default.pprof` -> Where to save the generated profile to.
4. Once finished, commit the newly generated profile file as `build/default.pprof` to your branch.
   It will be automatically picked up for the next release.
