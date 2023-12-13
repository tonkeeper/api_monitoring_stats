FROM golang:1.21 AS gobuild
WORKDIR /build-dir
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY config config
COPY services services
COPY *.go ./
RUN go build -v -o /tmp/api

FROM ubuntu:22.04 as monitoring
RUN mkdir -p /app/lib
RUN apt-get update && \
    apt-get install -y tzdata ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY --from=gobuild /tmp/api /app/api
CMD ["/app/api"]
