# FROM golang:1.24-bookworm AS base
#
# WORKDIR /app
#
# COPY go.mod go.sum ./
#
# RUN go mod download
#
# COPY . .
#
# RUN ["templ", "generate"]
# RUN go build -o mqtt-monitoring
#
# EXPOSE 8000
#
# CMD ["/build/mqtt-monitoring"]

# https://templ.guide/quick-start/installation/#docker
# Fetch
FROM golang:latest AS fetch-stage
COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download

# Generate
FROM ghcr.io/a-h/templ:latest AS generate-stage
COPY --chown=65532:65532 . /app
WORKDIR /app
RUN ["templ", "generate"]

FROM golang:latest AS deploy-stage
COPY --from=generate-stage /app /app
WORKDIR /app
RUN go build -o mqtt-monitoring -buildvcs=false

EXPOSE 8000

CMD ["/app/mqtt-monitoring"]
