FROM golang:1.20-alpine AS builder
RUN mkdir /build
WORKDIR /build
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk update && apk add git
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go install github.com/swaggo/swag/cmd/swag@f475da2  # Related to Issue #1564, waiting for the next release to solve this problem
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN swag init --exclude ./test
RUN go build -o /build/kayak .

FROM alpine:3.15
RUN mkdir /app
WORKDIR /app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk update
COPY --from=builder /build/kayak /app/kayak
COPY ./config.yaml /app/config.yaml
EXPOSE 9000
RUN mkdir ./log
CMD /app/kayak