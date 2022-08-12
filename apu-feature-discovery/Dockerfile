# global arguments for all build stages
ARG BASE_IMAGE_FULL
ARG VERSION

FROM golang:latest as build

WORKDIR /work

COPY ./ ./
RUN go build -o apu-feature-discovery .

# RUN go test .

FROM ${BASE_IMAGE_FULL} as full

COPY --from=build /work/apu-feature-discovery /usr/bin/apu-feature-discovery

LABEL io.k8s.display-name="LYNXI APU Feature Discovery Plugin"
LABEL name="LYNXI APU Feature Discovery Plugin"
LABEL vendor="LYNXI"
LABEL version="${VERSION}"
LABEL release="N/A"
LABEL summary="APU plugin to the node feature discovery for Kubernetes"
LABEL description="APU plugin to the node feature discovery for Kubernetes"
COPY ./LICENSE /licenses/LICENSE

ENTRYPOINT ["/usr/bin/apu-feature-discovery"]
