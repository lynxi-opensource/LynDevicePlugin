
FROM lynxidocker/lynxi-docker-ubuntu-18.04:1.4.0

ARG BIN

COPY ${BIN} main
CMD [ "./main" ]