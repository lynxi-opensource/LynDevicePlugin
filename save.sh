#!/bin/bash
set -e

version="1.9.0"
archs=(amd64 arm64)
targets=(lynxi-device-plugin lynxi-exporter apu-feature-discovery lynsmi-service lynxi-device-discovery)
image_prefix=192.168.9.41:5000

for target in ${targets[@]}; do
    for arch in ${archs[@]}; do
        docker pull ${image_prefix}/lynxidocker/${target}:${version} --platform ${arch}
        docker tag ${image_prefix}/lynxidocker/${target}:${version} lynxidocker/${target}:${version}
        docker save lynxidocker/${target}:${version} -o release/${target}-${version}-${arch}.tar
    done
done
