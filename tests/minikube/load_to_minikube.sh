#!/bin/bash
set -e

targets=(lynxi-device-plugin lynxi-exporter apu-feature-discovery lynsmi-service lynxi-device-discovery)

for target in ${targets[@]}; do
    minikube image load --overwrite=true release/${target}-1.9.0-${arch}.tar
done

minikube image load ubuntu:20.04
