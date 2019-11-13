FROM golang:latest

MAINTAINER zs "810909753@qq.com"

WORKDIR /app
ENV GO_SSO_WORKDIR  /app
ENV GOPROXY   https://gocenter.io

ADD . /app
RUN go build  -mod=vendor  main.go

ENTRYPOINT ["./main"]