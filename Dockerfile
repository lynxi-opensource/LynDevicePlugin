FROM golang:latest as builder

ENV LYNXI_VISIBLE_DEVICES=all
RUN go env -w GOPROXY=https://goproxy.cn,direct
WORKDIR /work
COPY go.mod go.mod
RUN go mod download
COPY ${BIN} ${BIN}
ARG BIN
RUN go build -o ${BIN} lyndeviceplugin/${BIN}

FROM ubuntu:18.04

WORKDIR /work
ARG BIN
COPY --from=builder /work/${BIN}/${BIN} main
CMD [ "./main" ]