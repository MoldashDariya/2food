FROM golang:1.25.4-alpine3.22 as build

ARG SERVICE_PATH
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/binary ${SERVICE_PATH}


FROM gcr.io/distroless/static-debian12:nonroot as backend

COPY --from=build --chown=nonroot:nonroot --chmod=755 /app/binary /server

USER nonroot

EXPOSE 8080

CMD ["/server"]


FROM gcr.io/distroless/static-debian12:nonroot as migrator

COPY --from=build --chown=nonroot:nonroot --chmod=755 /app/binary /migrator
COPY --chown=nonroot:nonroot ./migrations /migrations

USER nonroot

CMD ["/migrator"]
