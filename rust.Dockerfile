FROM --platform=$TARGETPLATFORM ubuntu:18.04

WORKDIR /work
ARG BIN
ARG TARGETARCH
COPY $TARGETARCH/lynsmi-service main
ENTRYPOINT [ "/work/main" ]