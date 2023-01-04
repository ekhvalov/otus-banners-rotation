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
ENV ROTATOR_RABBITMQ_HOST="rabbitmq"
ENV ROTATOR_RABBITMQ_PORT="5672"
ENV ROTATOR_RABBITMQ_USERNAME="guest"
ENV ROTATOR_RABBITMQ_PASSWORD="guest"
ENV ROTATOR_RABBITMQ_QUEUE_NAME="events"

EXPOSE 8081

COPY --from=build /opt/rotator /usr/local/bin/rotator
COPY ./configs/rotator.toml /etc/rotator/config.toml

CMD [ "rotator", "--config", "/etc/rotator/config.toml" ]
