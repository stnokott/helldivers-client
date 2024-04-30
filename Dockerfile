FROM alpine:3.19.1 AS tzbuild

RUN apk update && apk add --no-cache tzdata

FROM scratch

COPY --from=tzbuild ["/usr/share/zoneinfo", "/usr/share/zoneinfo"]
COPY ["scripts/migrations/*.sql", "scripts/migrations/"]
COPY ["app", "/app"]
ENTRYPOINT [ "/app" ]
