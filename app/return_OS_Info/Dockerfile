FROM golang:latest
WORKDIR /go/src/app
COPY ./app /go/src/app
RUN go get -u github.com/comail/colog
ENV PORT=8080
CMD ["go", "run", "returnJson.go"]
