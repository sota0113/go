FROM golang:latest
WORKDIR /go/src/app
RUN go get -u github.com/comail/colog
RUN go get -u github.com/peterbourgon/diskv
COPY ./app/returnJson.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest
WORKDIR /root
COPY --from=0 /go/src/app/app .
ENV PORT=8080
ENV MAIN_DIR="/list"
ENV HEALTH_DIR="/"
ENV LOG_TYPE="stdout"
CMD ["./app"]
