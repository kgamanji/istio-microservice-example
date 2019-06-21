FROM golang:1.11-alpine AS build_base
RUN apk add bash ca-certificates git gcc g++ libc-dev
WORKDIR /go/src/github.com/AlexsJones/golang-microservice-example

# Force the go compiler to use modules
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

# This image builds the weavaite server
FROM build_base AS server_builder
# Here we copy the rest of the source code
COPY . .
# And compile the project
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build && go install

FROM alpine as runner

COPY --from=server_builder /go/bin/golang-microservice-example /bin/gme
ENTRYPOINT ["/bin/gme"]
