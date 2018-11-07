FROM golang:1.9 AS build-env

MAINTAINER Johannes Amorosa

RUN mkdir -p /go/src/app
COPY . /go/src/app

WORKDIR /go/src/app

# download the dependencies and build the application
RUN go-wrapper download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go-wrapper install

FROM alpine:3.8 AS compress-env

COPY --from=build-env /go/bin/app /go/bin/app

RUN apk update && apk add upx && upx --best /go/bin/app

FROM alpine:3.8 AS runtime-env

COPY --from=compress-env /go/bin/app /usr/local/bin/app

EXPOSE 8081

WORKDIR /workdir
CMD ["mkdir /workdir/serve"]

ENTRYPOINT ["/usr/local/bin/app", "-port", "8081", "-pathPrefix", "/cool/fake/prefix/path", "-upload_limit", "1000000000", "-serverRoot", "/workdir/serve"]


