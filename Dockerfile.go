FROM golang:1.22 AS build

WORKDIR /app

COPY ./src/go /app

RUN go get github.com/andrewburto/mongodb-go && \
    go build -o mongo . && \
    chmod +x mongo

FROM busybox:latest

COPY --from=build /app/mongo /bin/mongo

CMD ["mongo"]