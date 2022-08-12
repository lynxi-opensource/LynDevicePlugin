
FROM lynxidocker/lynxi-docker-ubuntu-18.04:1.3.1

ARG BIN

COPY ${BIN} main
CMD [ "./main" ]