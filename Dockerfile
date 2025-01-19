FROM golang:alpine as builder

# ENV GO111MODULE=on

LABEL maintainer="Don Townsend <townofdon@gmail.com>"

# Install git
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod tidy
RUN go mod vendor

RUN go build

EXPOSE 8080

CMD ["./tutorial-go-rss-server"]