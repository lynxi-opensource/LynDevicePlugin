FROM registry.cn-hangzhou.aliyuncs.com/google_containers/kicbase:v0.0.44

ENV LYNXI_VISIBLE_DEVICES=all
COPY ./installer ./installer
RUN ./installer
COPY ./daemon.json /etc/docker/daemon.json