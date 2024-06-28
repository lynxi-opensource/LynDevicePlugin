#!/bin/bash
set -e

targets=(lynxi-device-plugin lynxi-exporter apu-feature-discovery lynsmi-service lynxi-device-discovery)

for target in ${targets[@]}; do
    docker load --input release/${target}-${version}-${arch}.tar
done
