
FROM golang:latest as builder

ENV LYNXI_VISIBLE_DEVICES=all
RUN go env -w GOPROXY=https://goproxy.cn,direct
WORKDIR /work
COPY go.mod go.mod
RUN go mod download
COPY . .
RUN go build -o lynxi-k8s-device-plugin lyndeviceplugin/lynxi-k8s-device-plugin
RUN go build -o lynxi-exporter lyndeviceplugin/lynxi-exporter

FROM ubuntu:18.04
ARG BIN

WORKDIR /work
COPY --from=builder /work/${BIN}/${BIN} main
CMD [ "./main" ]