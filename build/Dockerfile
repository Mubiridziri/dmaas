FROM golang:1.21.6

COPY .. app/
WORKDIR app
RUN go mod download
RUN go build ./cmd/dmaas

CMD ["./dmaas"]