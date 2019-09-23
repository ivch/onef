# Build app in a stock Go builder container
FROM golang:1.13-alpine3.10 as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ENV GO111MODULE=on

COPY . /go/src/onef
WORKDIR /go/src/onef

RUN cd cmd && go build -mod=vendor  -a -o /go/bin/svc

## Pull binaries into a second stage deploy alpine container
FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/svc /svc/

WORKDIR /svc

RUN chmod +x svc

CMD ["./svc"]