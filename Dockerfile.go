FROM golang:1.22 AS build

WORKDIR /app

COPY ./src/go /app

RUN go get . && \
    go build -o mongo . && \
    chmod +x mongo

FROM busybox:latest

EXPOSE 8080

COPY --from=build /app/mongo /bin/mongo

CMD ["mongo"]