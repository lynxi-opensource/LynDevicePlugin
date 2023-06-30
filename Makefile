version = 1.5.0

build:
	./build.sh

chart:
	mkdir release -p
	cd release && helm package ../LynDevicePlugin

namespace:
	kubectl create namespace device-plugin

example:
	kubectl apply -f example.yml

install:
	helm install -n device-plugin --set imagePullPolicy=Always lynxi-device-plugin release/LynDevicePlugin-${version}.tgz

install-no-service-monitor:
	helm install -n device-plugin --set imagePullPolicy=Always --set lynxiExporter.serviceMonitor.enable=false lynxi-device-plugin release/LynDevicePlugin-${version}.tgz

upgrade:
	helm upgrade -n device-plugin lynxi-device-plugin release/LynDevicePlugin-${version}.tgz

uninstall:
	helm uninstall -n device-plugin lynxi-device-plugin

list:
	helm list -n device-plugin