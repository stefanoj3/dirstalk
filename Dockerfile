FROM golang:1.15-alpine3.12 as builder

RUN adduser -D -g '' dirstalkuser
RUN apk add --update make git ca-certificates

RUN mkdir -p $GOPATH/src/github.com/stefanoj3/dirstalk
ADD . $GOPATH/src/github.com/stefanoj3/dirstalk
WORKDIR $GOPATH/src/github.com/stefanoj3/dirstalk

RUN make dep
RUN make build

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/github.com/stefanoj3/dirstalk/dist/dirstalk /bin/dirstalk
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER dirstalkuser
CMD ["dirstalk"]