FROM lynxidocker/lynxi-docker-ubuntu-18.04:1.3.1 as builder

ENV LYNXI_VISIBLE_DEVICES=all
RUN apt update
RUN apt install git gcc g++ -y
RUN apt install wget -y
RUN wget https://golang.google.cn/dl/go1.18.1.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.1.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w GOFLAGS=-buildvcs=false
WORKDIR /work
COPY go.mod go.mod
RUN go mod download