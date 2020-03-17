FROM golang:1.14
WORKDIR Ampel
COPY go.* .
RUN go mod download
COPY *.go .
RUN CGO_ENABLED=0 go build

FROM eu.gcr.io/vseth-public/base:delta
COPY --from=0 /go/Ampel/Ampel .
COPY src src
EXPOSE 8080