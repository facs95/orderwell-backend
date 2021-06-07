FROM golang:latest

WORKDIR /app

COPY . /app

RUN go install

CMD ["/go/bin/orderwell-backend"]
