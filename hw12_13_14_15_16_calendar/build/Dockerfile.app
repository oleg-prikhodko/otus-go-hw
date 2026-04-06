FROM golang:1.23 AS build
ARG SERVICE_NAME
ARG LDFLAGS

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o ${SERVICE_NAME} ./cmd/${SERVICE_NAME}

FROM alpine:3.9
ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

RUN apk add --no-cache ca-certificates bash gettext

WORKDIR /app

COPY --from=build /build/${SERVICE_NAME} /app/${SERVICE_NAME}
COPY deployments/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD /app/${SERVICE_NAME} -config /tmp/config.yaml
