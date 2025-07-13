# https://templ.guide/quick-start/installation/#docker
# Fetch
FROM golang:latest AS fetch-stage
COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download


FROM node:20 AS frontend-build-stage
WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm install

RUN mkdir -p /app/frontend/static/css \
  && cp node_modules/@picocss/pico/css/pico.jade.min.css /app/frontend/static/css/


FROM ghcr.io/a-h/templ:latest AS generate-stage
COPY --chown=65532:65532 . /app
WORKDIR /app
RUN ["templ", "generate"]


FROM golang:latest AS deploy-stage
COPY --from=generate-stage /app /app
COPY --from=frontend-build-stage /app/frontend/static /app/frontend/static

WORKDIR /app
RUN go build -o mqtt-monitoring -buildvcs

EXPOSE 8000

CMD ["/app/mqtt-monitoring"]
