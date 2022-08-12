REPO=lynxidocker
VERSION=0.2.1
IMG=$(REPO)/lynxi-k8s-device-plugin:$(VERSION)

build: 
	docker build -t $(IMG) .

push: build 
	docker push $(IMG)

apply: push 
	kubectl apply -f DaemonSet.yml

delete:  
	kubectl delete -f DaemonSet.yml

replace: push 
	kubectl replace -f DaemonSet.yml

force_replace: push 
	kubectl replace -f DaemonSet.yml --force

MAIN_PKG=lyndeviceplugin/lynxi-k8s-device-plugin
PKGS=$(MAIN_PKG) $(MAIN_PKG)/allocator $(MAIN_PKG)/server $(MAIN_PKG)/service

test:
	go test -run=. -timeout=10s -race -v $(PKGS)

coverage:
	go test -run=. -v -covermode=count -coverprofile=coverage.out $(PKGS)
	go tool cover -func=coverage.out
	rm coverage.out

docs:
	godocdown

run:
	# TODO

run_local:
	# TODO