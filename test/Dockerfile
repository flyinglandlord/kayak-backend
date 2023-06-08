FROM golang:1.18-alpine
RUN mkdir /app
WORKDIR /app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk update && apk add git gcc musl-dev
RUN go env -w GOPROXY=https://goproxy.cn,direct
COPY . .
RUN go mod download
CMD go test -v -coverpkg=./... -coverprofile=profile.cov ./test && go tool cover -func profile.cov
