lynxi-k8s-device-plugin-version = 0.2.1
lynxi-exporter-version = 0.2.1

build_with_docker:
	docker build -t lyndeviceplugin_image_for_build:latest . -f build.Dockerfile
	docker run --rm -e LYNXI_VISIBLE_DEVICES=all -v $(PWD):/work lyndeviceplugin_image_for_build:latest bash ./build.sh

push-lynxi-k8s-device-plugin: build_with_docker
	docker build -t lynxidocker/lynxi-k8s-device-plugin:${lynxi-k8s-device-plugin-version} release -f Dockerfile --build-arg BIN=lynxi-k8s-device-plugin
	docker push lynxidocker/lynxi-k8s-device-plugin:${lynxi-k8s-device-plugin-version}

push-lynxi-exporter: build_with_docker
	docker build -t lynxidocker/lynxi-exporter:${lynxi-exporter-version} release -f Dockerfile --build-arg BIN=lynxi-exporter
	docker push lynxidocker/lynxi-exporter:${lynxi-exporter-version}

build: 
	./build.sh

chart:
	mkdir release -p
	cd release && helm package ../lynxi-device-chart