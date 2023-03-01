lynxi-k8s-device-plugin-version = 1.2.0
lynxi-exporter-version = 1.2.0 

DEVICE-PLUGIN-IMAGE = lynxidocker/lynxi-k8s-device-plugin:${lynxi-k8s-device-plugin-version}
EXPORTER-IMAGE = lynxidocker/lynxi-exporter:${lynxi-exporter-version}

build-amd64:
	docker build -t ${DEVICE-PLUGIN-IMAGE}-amd64 . -f Dockerfile --build-arg BIN=lynxi-k8s-device-plugin
	docker build -t ${EXPORTER-IMAGE}-amd64 . -f Dockerfile --build-arg BIN=lynxi-exporter

push-amd64:
	docker push ${DEVICE-PLUGIN-IMAGE}-amd64
	docker push ${EXPORTER-IMAGE}-amd64

build-arm64:
	docker build -t ${DEVICE-PLUGIN-IMAGE}-arm64 . -f Dockerfile --build-arg BIN=lynxi-k8s-device-plugin
	docker build -t ${EXPORTER-IMAGE}-arm64 . -f Dockerfile --build-arg BIN=lynxi-exporter

push-arm64:
	docker push ${DEVICE-PLUGIN-IMAGE}-arm64
	docker push ${EXPORTER-IMAGE}-arm64

docker-manifest:
	docker manifest create ${DEVICE-PLUGIN-IMAGE} ${DEVICE-PLUGIN-IMAGE}-amd64 ${DEVICE-PLUGIN-IMAGE}-arm64
	docker manifest annotate ${DEVICE-PLUGIN-IMAGE} ${DEVICE-PLUGIN-IMAGE}-amd64 --os linux --arch amd64
	docker manifest annotate ${DEVICE-PLUGIN-IMAGE} ${DEVICE-PLUGIN-IMAGE}-arm64 --os linux --arch arm64
	docker manifest push ${DEVICE-PLUGIN-IMAGE}

	docker manifest create ${EXPORTER-IMAGE} ${EXPORTER-IMAGE}-amd64 ${EXPORTER-IMAGE}-arm64
	docker manifest annotate ${EXPORTER-IMAGE} ${EXPORTER-IMAGE}-amd64 --os linux --arch amd64
	docker manifest annotate ${EXPORTER-IMAGE} ${EXPORTER-IMAGE}-arm64 --os linux --arch arm64
	docker manifest push ${EXPORTER-IMAGE}


chart:
	mkdir release -p
	cd release && helm package ../LynDevicePlugin

install:
	helm install -n device-plugin lynxi-device-plugin release/LynDevicePlugin-${lynxi-k8s-device-plugin-version}.tgz

upgrade:
	helm upgrade -n device-plugin lynxi-device-plugin release/LynDevicePlugin-${lynxi-k8s-device-plugin-version}.tgz

uninstall:
	helm uninstall -n device-plugin lynxi-device-plugin

list:
	helm list -n device-plugin