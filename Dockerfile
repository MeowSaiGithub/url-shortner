# syntax=docker/dockerfile:1
##
## Build
##

FROM golang:1.20.5-alpine3.17 AS build

WORKDIR /app

RUN apk --no-cache add tzdata

COPY internal ./internal
COPY go.mod ./go.mod
COPY go.sum ./go.sum
COPY *.go ./

RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o url-shortner -ldflags '-w -s' .

##
## Deploy
##
#FROM alpine:3.17
#
#WORKDIR /app
#
#COPY --from=build /app/url-shortner ./url-shortner
#
#RUN chmod +x ./url-shortner
#
#CMD ./url-shortner

##
##  Smaller Deploy
##
FROM scratch

WORKDIR /app

COPY --from=build /app/url-shortner ./url-shortner
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
#COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app/url-shortner"]