FROM golang:1.19-alpine

WORKDIR /app/loadbalancer

COPY go.mod ./

RUN go mod download

COPY ./ .

RUN go build -o loadbalancer .

EXPOSE 8000

CMD ["./loadbalancer"]