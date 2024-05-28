#!/bin/bash
set -e

version="1.7.1"
archs=(amd64 arm64)
go_targets=(lynxi-device-plugin lynxi-exporter apu-feature-discovery)
out_dir=bin

export CGO_ENABLED=0
for target in ${go_targets[@]}; do
    echo "building ${target}"
    for arch in ${archs[@]}; do
        GOOS=linux GOARCH=${arch} go build -o ${out_dir}/${arch}/${target} lyndeviceplugin/${target}
    done
    docker buildx build --platform linux/amd64,linux/arm64 -t lynxidocker/${target}:${version} ${out_dir} -f go.Dockerfile --build-arg BIN=${target} --push
done

echo "building lynsmi-service"
cargo build -r --target=x86_64-unknown-linux-gnu -p lynsmi-service
cp target/x86_64-unknown-linux-gnu/release/lynsmi-service ${out_dir}/amd64/
cargo build -r --target=aarch64-unknown-linux-gnu -p lynsmi-service
cp target/aarch64-unknown-linux-gnu/release/lynsmi-service ${out_dir}/arm64/
docker buildx build --platform linux/amd64,linux/arm64 -t lynxidocker/lynsmi-service:${version} ${out_dir} -f rust.Dockerfile --build-arg BIN=lynsmi-service --push

echo "building lynxi-device-discovery"
cargo build -r --target=x86_64-unknown-linux-gnu -p lynxi-device-discovery
cp target/x86_64-unknown-linux-gnu/release/lynxi-device-discovery ${out_dir}/amd64/
cargo build -r --target=aarch64-unknown-linux-gnu -p lynxi-device-discovery
cp target/aarch64-unknown-linux-gnu/release/lynxi-device-discovery ${out_dir}/arm64/
docker buildx build --platform linux/amd64,linux/arm64 -t lynxidocker/lynxi-device-discovery:${version} ${out_dir} -f rust.Dockerfile --build-arg BIN=lynxi-device-discovery --push
