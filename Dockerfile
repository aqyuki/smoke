FROM golang:1.24.2 AS build

WORKDIR /app
RUN --mount=type=bind,target=. go mod download
RUN --mount=type=bind,target=. go mod verify
RUN --mount=type=bind,target=. go build -o /dist/smoke main.go

FROM gcr.io/distroless/cc-debian12 AS prod

ENV TZ=Asia/Tokyo

WORKDIR /app
COPY --from=build --chown=root:root /dist/smoke /app/smoke
STOPSIGNAL SIGINT
ENTRYPOINT ["./smoke"]
