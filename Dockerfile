ARG GO_VERSION=1.9.3
ARG ALPINE_VERSION=3.7

FROM golang:${GO_VERSION} AS BUILD
LABEL maintainer="frankgreco@northwesternmutual.com"
LABEL version="${VERSION}"
ARG VERSION=""
WORKDIR /go/src/github.com/northwesternmutual/kanali/
COPY Gopkg.toml Gopkg.lock Makefile /go/src/github.com/northwesternmutual/kanali/
RUN make install
COPY ./ /go/src/github.com/northwesternmutual/kanali/
RUN CGO_ENABLED=0 make binary

FROM alpine:${ALPINE_VERSION}
LABEL maintainer="frankgreco@northwesternmutual.com"
LABEL version="${VERSION}"
COPY --from=BUILD /go/bin/kanali /
COPY --from=BUILD /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/kanali"]
