FROM lynxidocker/lynxi-docker-ubuntu-18.04:1.3.1 as builder

RUN apt update
RUN apt install git gcc g++ -y
RUN apt install wget -y
RUN wget https://golang.google.cn/dl/go1.18.1.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.1.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
WORKDIR /work
RUN go env -w GOPROXY=https://goproxy.cn,direct
COPY go.mod go.mod
RUN go mod download
COPY ./ ./
ENV LYNXI_VISIBLE_DEVICES=all
RUN go build -o lynxi-exporter .

FROM lynxidocker/lynxi-docker-ubuntu-18.04:1.3.1

COPY --from=builder /work/lynxi-exporter /usr/bin/lynxi-exporter
CMD [ "lynxi-exporter" ]
