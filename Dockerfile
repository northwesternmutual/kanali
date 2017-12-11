ARG GO_VERSION=1.8.5

FROM golang:${GO_VERSION}-jessie AS BUILD
LABEL maintainer="frankgreco@northwesternmutual.com"
LABEL version="${VERSION}"
ARG VERSION=""
ARG GLIDE_VERSION=0.12.3
WORKDIR /go/src/github.com/northwesternmutual/kanali/
RUN wget "https://github.com/Masterminds/glide/releases/download/v${GLIDE_VERSION}/glide-v${GLIDE_VERSION}-`go env GOHOSTOS`-`go env GOHOSTARCH`.tar.gz" -O /tmp/glide.tar.gz \
    && mkdir /tmp/glide \
    && tar --directory=/tmp/glide -xvf /tmp/glide.tar.gz \
    && rm -rf /tmp/glide.tar.gz \
    && export PATH=$PATH:/tmp/glide/`go env GOHOSTOS`-`go env GOHOSTARCH`
COPY glide.lock glide.yaml Makefile /go/src/github.com/northwesternmutual/kanali/
RUN make install
COPY ./ /go/src/github.com/northwesternmutual/kanali/
RUN GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` make binary ${VERSION}

FROM debian:jessie-slim
LABEL maintainer="frankgreco@northwesternmutual.com"
LABEL version="${VERSION}"
COPY --from=BUILD /go/bin/kanali /
ENTRYPOINT ["/kanali"]