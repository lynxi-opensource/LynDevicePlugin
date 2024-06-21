#!/bin/bash
set -e

version="1.9.0"
arch=arm64
targets=(lynxi-device-plugin lynxi-exporter apu-feature-discovery lynsmi-service lynxi-device-discovery)

for target in ${targets[@]}; do
    docker load --input release/${target}-${version}-${arch}.tar
done
