version = 1.10.0
image_prefix = 192.168.9.41:5000/
arch = amd64

build:
	version=${version} image_prefix=${image_prefix} ./scripts/build.sh

save:
	version=${version} image_prefix=${image_prefix} ./scripts/save.sh

load_amd64:
	version=${version} arch=${arch} ./scripts/load_amd64.sh

minikube_load_images:
	version=${version} arch=${arch} .tests/minikube/load_to_minikube.sh

chart:
	mkdir release -p
	cd release && helm package --app-version ${version} --version ${version} ../LynDevicePlugin

example:
	kubectl apply -f scripts/example.yml

example-uninstall:
	kubectl delete -f scripts/example.yml

install:
	helm install --set imagePullPolicy=Never ldp release/LynDevicePlugin-${version}.tgz

install-no-service-monitor:
	helm install --set imagePullPolicy=Never --set lynxiExporter.serviceMonitor.enable=false --set namespace.name=default --set namespace.create=false ldp release/LynDevicePlugin-${version}.tgz

upgrade:
	helm upgrade ldp release/LynDevicePlugin-${version}.tgz

uninstall:
	helm uninstall ldp
