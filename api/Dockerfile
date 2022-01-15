FROM golang:1.17.6

WORKDIR /go/src/app

# COPY *.go ./
COPY ./api ./go/src/app

RUN go mod init example.com/app
RUN go mod tidy

EXPOSE 80