FROM scratch

COPY ["scripts/migrations/*.sql", "scripts/migrations/"]
COPY ["app", "/app"]
ENTRYPOINT [ "/app" ]
