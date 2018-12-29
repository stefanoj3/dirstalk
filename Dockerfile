FROM golang:1.11.4-alpine3.8 as builder

RUN adduser -D -g '' dirstalkuser
RUN apk add --update make git

RUN mkdir -p $GOPATH/src/github.com/stefanoj3/dirstalk
ADD . $GOPATH/src/github.com/stefanoj3/dirstalk
WORKDIR $GOPATH/src/github.com/stefanoj3/dirstalk

RUN make dep
RUN make build

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/github.com/stefanoj3/dirstalk/dist/linux_build /bin/dirstalk

USER dirstalkuser
CMD ["dirstalk"]