.PHONY: build push

REPO=lynxidocker
VERSION=0.2.1
IMG=$(REPO)/lynxi-exporter:$(VERSION)

build: 
	docker build -t $(IMG) .

push: build 
	docker push $(IMG)

run: 
	go run .

curl: 
	curl http://localhost:2112/metrics | grep lynxi

daemonset: push 
	kubectl replace -f deploy/DaemonSet.yml 

force_replace_daemonset: push 
	kubectl replace -f deploy/DaemonSet.yml --force

service: daemonset 
	kubectl replace -f deploy/Service.yml 

service_monitor: service 
	kubectl replace -f deploy/ServiceMonitor.yml