FROM golang:1.16-alpine as local-dev
RUN apk add --no-cache git curl build-base autoconf automake libtool make bash protobuf mysql-client \
    && GO111MODULE=off go get -u google.golang.org/protobuf/cmd/protoc-gen-go \
    && GO111MODULE=off go get -u github.com/google/wire/cmd/wire \
    && go get github.com/twitchtv/twirp/protoc-gen-twirp@v8.0.0
ENV TZ=Asia/Tokyo


FROM golang:1.16-alpine as builder
ENV APP_ROOT /app
WORKDIR $APP_ROOT
ARG _ENTRYPOINT
COPY . $APP_ROOT/
RUN apk add --no-cache git curl build-base autoconf automake libtool make bash protobuf \
    && GO111MODULE=off go get -u google.golang.org/protobuf/cmd/protoc-gen-go \
    && GO111MODULE=off go get -u github.com/google/wire/cmd/wire \
    && go get github.com/twitchtv/twirp/protoc-gen-twirp@v8.0.0
RUN make gen && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo $APP_ROOT/$_ENTRYPOINT


FROM alpine as deploy
ENV APP_ROOT /app
WORKDIR $APP_ROOT
RUN apk add ca-certificates
COPY --from=builder $APP_ROOT/main $APP_ROOT/main
ENV TZ=Asia/Tokyo