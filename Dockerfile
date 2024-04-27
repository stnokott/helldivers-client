FROM scratch

# TODO: get TZ data to work, see e.g. https://stackoverflow.com/questions/47794740/gos-time-doesnt-work-under-the-docker-image-from-scratch

COPY ["scripts/migrations/*.sql", "scripts/migrations/"]
COPY ["app", "/app"]
ENTRYPOINT [ "/app" ]
