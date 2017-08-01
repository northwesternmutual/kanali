# multi-stage build
# build stage
# what version of the go language should be used
ARG GO_VERSION=1.8.1
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
# copy all of our files into the container
COPY ./ /go/src/github.com/northwesternmutual/kanali/
# Install Glide - our go dependency management tool
RUN wget "https://github.com/Masterminds/glide/releases/download/v${GLIDE_VERSION}/glide-v${GLIDE_VERSION}-`go env GOHOSTOS`-`go env GOHOSTARCH`.tar.gz" -O /tmp/glide.tar.gz \
    && mkdir /tmp/glide \
    && tar --directory=/tmp/glide -xvf /tmp/glide.tar.gz \
    && rm -rf /tmp/glide.tar.gz
# Install dependencies
RUN export PATH=$PATH:/tmp/glide/`go env GOHOSTOS`-`go env GOHOSTARCH` \
    && make install
# Set project version
RUN sed -ie "s/changeme/`echo ${VERSION}`/g" /go/src/github.com/northwesternmutual/kanali/cmd/version.go
# Build project
RUN make build
# add ca certificates bundle
RUN curl http://curl.haxx.se/ca/cacert.pem -o /etc/pki/tls/certs/ca-bundle.crt
# set our entrypoint
ENTRYPOINT ["./kanali"]