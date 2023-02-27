lynxi-k8s-device-plugin-version = 1.2.0
lynxi-exporter-version = 1.2.0 


build:
	docker build -t lynxidocker/lynxi-k8s-device-plugin:${lynxi-k8s-device-plugin-version} . -f Dockerfile --build-arg BIN=lynxi-k8s-device-plugin
	docker build -t lynxidocker/lynxi-exporter:${lynxi-exporter-version} . -f Dockerfile --build-arg BIN=lynxi-exporter

push: build
	docker push lynxidocker/lynxi-k8s-device-plugin:${lynxi-k8s-device-plugin-version}
	docker push lynxidocker/lynxi-exporter:${lynxi-exporter-version}

hs110-kubeedge-build:
	docker build -t lynxidocker/lynxi-k8s-device-plugin-hs110-kubeedge:${lynxi-k8s-device-plugin-version} . -f Dockerfile --build-arg BIN=lynxi-k8s-device-plugin
	docker build -t lynxidocker/lynxi-exporter-hs110-kubeedge:${lynxi-exporter-version} . -f Dockerfile --build-arg BIN=lynxi-exporter

hs110-kubeedge-push: hm100-build
	docker push lynxidocker/lynxi-k8s-device-plugin-hs110-kubeedge:${lynxi-k8s-device-plugin-version}
	docker push lynxidocker/lynxi-exporter-hs110-kubeedge:${lynxi-exporter-version}

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