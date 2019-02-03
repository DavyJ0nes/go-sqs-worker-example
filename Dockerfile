# Multistage Build
# -- Builder Image
FROM golang:1.11.4-alpine3.8 As Builder

RUN apk update && \
    apk add curl ca-certificates git && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/davyj0nes/worker-example

COPY ./Gopkg.toml Gopkg.toml
COPY ./Gopkg.lock Gopkg.lock
RUN dep ensure -vendor-only

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo --installsuffix netgo -o worker .

# Main Image
FROM alpine:3.8 AS alpine
RUN adduser -D -u 10001 dockmaster
RUN apk update && \
    apk add ca-certificates

# Copy binary from builder image
COPY --from=Builder /go/src/github.com/davyj0nes/worker-example/worker /bin/worker
RUN chmod a+x /bin/worker

USER dockmaster

EXPOSE 8080
CMD ["worker"]
