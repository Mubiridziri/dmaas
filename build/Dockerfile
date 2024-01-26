FROM golang:1.21.6

COPY .. app/
WORKDIR app
RUN go mod download
RUN go build -o /app-build

CMD ["/app-build", "-release-mode=true"]