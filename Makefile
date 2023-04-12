version = 1.4.0
targets = lynxi-device-plugin lynxi-exporter

build-amd64:
	@for target in $(targets); do \
        docker build -t lynxidocker/$$target:$(version)-amd64 . -f Dockerfile --build-arg BIN=$$target; \
    done
	cd apu-feature-discovery && make build-amd64

push-amd64:
	@for target in $(targets); do \
		docker push lynxidocker/$$target:$(version)-amd64; \
    done
	cd apu-feature-discovery && make push-amd64

build-arm64:
	@for target in $(targets); do \
        docker build -t lynxidocker/$$target:$(version)-arm64 . -f Dockerfile --build-arg BIN=$$target; \
    done
	cd apu-feature-discovery && make build-arm64

push-arm64:
	@for target in $(targets); do \
		docker push lynxidocker/$$target:$(version)-arm64; \
    done
	cd apu-feature-discovery && make push-arm64

docker-manifest:
	@for target in $(targets); do \
		image = lynxidocker/$$target:$(version); \
		docker manifest create ${image} ${image}-amd64 ${image}-arm64; \
		docker manifest annotate ${image} ${image}-amd64 --os linux --arch amd64; \
		docker manifest annotate ${image} ${image}-arm64 --os linux --arch arm64; \
		docker manifest push ${image}; \
    done
	cd apu-feature-discovery && make docker-manifest


chart:
	mkdir release -p
	cd release && helm package ../LynDevicePlugin

namespace:
	kubectl create namespace device-plugin

example:
	kubectl apply -f example.yml

install:
	helm install -n device-plugin lynxi-device-plugin release/LynDevicePlugin-${version}.tgz

install-no-service-monitor:
	helm install -n device-plugin --set lynxiExporterServiceMonitor.enable=false lynxi-device-plugin release/LynDevicePlugin-${version}.tgz

upgrade:
	helm upgrade -n device-plugin lynxi-device-plugin release/LynDevicePlugin-${version}.tgz

uninstall:
	helm uninstall -n device-plugin lynxi-device-plugin

list:
	helm list -n device-plugin