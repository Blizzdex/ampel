#Stage 0 used to compile the go code
FROM golang:1.14
WORKDIR Ampel
COPY go.* .
RUN go mod download
COPY *.go .
RUN CGO_ENABLED=0 go build

# Stage 1 is based on the vis base image.
FROM eu.gcr.io/vseth-public/base:delta
COPY --from=0 /go/Ampel/Ampel .
COPY src src
EXPOSE 8080