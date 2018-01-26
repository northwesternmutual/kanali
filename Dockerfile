ARG GO_VERSION=1.9.2
ARG CENTOS_VERSION=7

FROM golang:${GO_VERSION} AS BUILD
LABEL maintainer="frankgreco@northwesternmutual.com"
LABEL version="${VERSION}"
ARG VERSION=""
WORKDIR /go/src/github.com/northwesternmutual/kanali/
COPY Gopkg.toml Gopkg.lock Makefile /go/src/github.com/northwesternmutual/kanali/
RUN make install

RUN cp -R vendor/* /go/src/ && rm -rf vendor

RUN cd /go/src/github.com/northwesternmutual/ && \
    git clone https://github.com/northwesternmutual/kanali-plugin-apikey.git && \
    cd kanali-plugin-apikey && \
    git checkout etcd-grpc && \
    make install && \
    cp -R vendor/* /go/src/ && \
    rm -rf vendor

COPY ./ /go/src/github.com/northwesternmutual/kanali/
RUN make binary
RUN cd /go/src/github.com/northwesternmutual/kanali-plugin-apikey && \
    go build -buildmode=plugin -o apiKey_v2.0.0-rc.1.so

FROM centos:${CENTOS_VERSION}
LABEL maintainer="frankgreco@northwesternmutual.com"
LABEL version="${VERSION}"
COPY --from=BUILD /go/bin/kanali /
COPY --from=BUILD /go/src/github.com/northwesternmutual/kanali-plugin-apikey/apiKey_v2.0.0-rc.1.so /plugins/
ENTRYPOINT ["/kanali"]