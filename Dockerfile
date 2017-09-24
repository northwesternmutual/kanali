ARG GO_VERSION=1.8.3
ARG CENTOS_VERSION=7

FROM golang:${GO_VERSION} AS BUILD
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
RUN sed -ie "s/changeme/`echo ${VERSION}`/g" /go/src/github.com/northwesternmutual/kanali/cmd/version.go
RUN curl -O https://raw.githubusercontent.com/northwesternmutual/kanali-plugin-apikey/v1.2.0/plugin.go
RUN GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -buildmode=plugin -o apiKey_v1.2.0.so plugin.go
RUN GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o kanali

FROM centos:${CENTOS_VERSION}
LABEL maintainer="frankgreco@northwesternmutual.com"
LABEL version="${VERSION}"
COPY --from=BUILD /go/src/github.com/northwesternmutual/kanali/apiKey_v1.2.0.so  /go/src/github.com/northwesternmutual/kanali/kanali /
RUN cp /apiKey_v1.2.0.so /apiKey.so
ENTRYPOINT ["/kanali"]