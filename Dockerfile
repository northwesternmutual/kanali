# multi-stage build
# build stage
# what version of the go language should be used
ARG GO_VERSION=1.8.3
# base build stage image from go version
FROM golang:${GO_VERSION} AS build-stage
# set image maintainer
MAINTAINER frankgreco@northwesternmutual.com
# what version of our project are we building
ARG VERSION="unknown version"
# what version of Glide should we use
ARG GLIDE_VERSION=0.12.3
# set our working directory
WORKDIR /go/src/github.com/northwesternmutual/kanali/
# Install Glide - our go dependency management tool
RUN wget "https://github.com/Masterminds/glide/releases/download/v${GLIDE_VERSION}/glide-v${GLIDE_VERSION}-`go env GOHOSTOS`-`go env GOHOSTARCH`.tar.gz" -O /tmp/glide.tar.gz \
    && mkdir /tmp/glide \
    && tar --directory=/tmp/glide -xvf /tmp/glide.tar.gz \
    && rm -rf /tmp/glide.tar.gz
RUN export PATH=$PATH:/tmp/glide/`go env GOHOSTOS`-`go env GOHOSTARCH`
# copy file necessary for dependency installation
COPY glide.lock glide.yaml Makefile /go/src/github.com/northwesternmutual/kanali/
# Install dependencies
RUN make install
# copy rest of source code into container
COPY ./ /go/src/github.com/northwesternmutual/kanali/
# Set project version
RUN sed -ie "s/changeme/`echo ${VERSION}`/g" /go/src/github.com/northwesternmutual/kanali/cmd/version.go
# download apikey plugin
RUN curl -O https://raw.githubusercontent.com/northwesternmutual/kanali-plugin-apikey/master/plugin.go
# compile plugin
RUN GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -buildmode=plugin -o apiKey.so plugin.go
# Build project
RUN GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o kanali

# production stage
FROM centos:latest
# set maintainer
MAINTAINER frankgreco@northwesternmutual.com
# add ca certificates bundle
RUN curl http://curl.haxx.se/ca/cacert.pem -o /etc/pki/tls/certs/ca-bundle.crt
# load plugin
COPY --from=build-stage /go/src/github.com/northwesternmutual/kanali/apiKey.so ./apiKey.so
# copy compiled binary from our build stage
COPY --from=build-stage /go/src/github.com/northwesternmutual/kanali/kanali .
# set our entrypoint
ENTRYPOINT ["/kanali"]