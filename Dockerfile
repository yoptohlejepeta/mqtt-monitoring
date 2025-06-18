FROM golang:1.24-bookworm AS base

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o mqtt-monitoring

EXPOSE 8000

CMD ["/build/mqtt-monitoring"]
