FROM golang:1.24 AS gobuild
WORKDIR /build-dir
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY config config
COPY services services
COPY *.go ./
RUN go build -v -o /tmp/api_monitoring_stats

FROM ubuntu:22.04 as monitoring
RUN mkdir -p /app/lib
RUN apt-get update && \
    apt-get install -y tzdata ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY --from=gobuild /tmp/api_monitoring_stats /app/api_monitoring_stats
CMD ["/app/api_monitoring_stats"]
