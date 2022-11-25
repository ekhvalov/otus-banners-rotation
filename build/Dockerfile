FROM golang:1.19.3 as build

ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 go build \
        -o /opt/rotator cmd/rotator/*

FROM alpine:3.16

ENV ROTATOR_GRPC_HOST="0.0.0.0"
ENV ROTATOR_GRPC_PORT="8081"
ENV ROTATOR_REDIS_HOST="redis"
ENV ROTATOR_REDIS_PORT="6379"
EXPOSE 8081

COPY --from=build /opt/rotator /usr/local/bin/rotator
COPY ./configs/rotator.toml /etc/rotator/config.toml

CMD [ "rotator", "--config", "/etc/rotator/config.toml" ]
