FROM golang:1.22

WORKDIR /app

COPY ./src/go /app

RUN go get github.com/andrewburto/mongodb-go && \
    go build -o mongo . && \
    chmod +x mongo

CMD ["./mongo"]