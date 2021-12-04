ARG goversion=1.15

# Stage 1 used to compile the go code
FROM golang:${goversion} as server
WORKDIR ampel2
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 go build

# Stage 1 is based on the vis base image.
FROM eu.gcr.io/vseth-public/base:delta
COPY --from=server /go/ampel2/ampel2 .
COPY src src
COPY migrations migrations
COPY cinit.yml /etc/cinit.d/ampel2.yml
