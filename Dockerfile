#Stage 0 used to compile the go code
ARG goversion=1.14

FROM golang:${goversion}
WORKDIR ampel2
COPY go.* ./
RUN go mod download
COPY *.go ./
COPY ampel ampel
RUN CGO_ENABLED=0 go build

# Stage 1 is based on the vis base image.
FROM eu.gcr.io/vseth-public/base:delta
COPY --from=0 /go/ampel2/ampel2 .
COPY src src
COPY migrations migrations
COPY cinit.yml /etc/cinit.d/ampel2.yml