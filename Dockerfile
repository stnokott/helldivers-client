FROM scratch

COPY ["migrations/*.json", "migrations/"]
COPY ["app", "/app"]
ENTRYPOINT [ "/app" ]
