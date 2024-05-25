FROM alpine:3.19.1

RUN apk update && \
	apk add --no-cache tzdata ca-certificates && \
	update-ca-certificates && \
	rm -rf /var/cache/apk/*

COPY ["scripts/migrations/*.sql", "scripts/migrations/"]
COPY ["app", "/app"]
ENTRYPOINT [ "/app" ]
