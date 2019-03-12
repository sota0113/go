FROM golang:latest
WORKDIR /go/src/app
COPY ./app /go/src/app
RUN go get -u github.com/comail/colog
RUN go get -u github.com/peterbourgon/diskv
ENV PORT=8080
ENV MAIN_DIR="/list"
ENV HEALTH_DIR="/"
ENV LOG_TYPE="stdout"
CMD ["go", "run", "returnJson.go"]
