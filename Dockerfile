FROM debian:jessie-slim
LABEL maintainer="frankgreco@northwesternmutual.com"
COPY kanali /
ENTRYPOINT ["/kanali"]