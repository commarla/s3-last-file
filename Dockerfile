# build stage
FROM golang:buster AS build-env
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o s3-last-object-time

# final stage
FROM debian:buster
COPY --from=build-env /app/s3-last-object-time /app/
RUN apt-get update && apt-get install -yq ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*
ENTRYPOINT ["/app/s3-last-object-time"]
