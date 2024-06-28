set -e

docker pull registry.cn-hangzhou.aliyuncs.com/google_containers/kicbase:v0.0.44
docker build -t minikube-base:v0.0.44 ./test_utils/minikube/
minikube delete
minikube start --driver docker --container-runtime docker --base-image=minikube-base:v0.0.44 --feature-gates=DynamicResourceAllocation=true --extra-config=apiserver.runtime-config=resource.k8s.io/v1alpha2=true --addons dashboard --addons metrics-server --docker-opt=add-runtime=lynxi=/usr/bin/lynxi-container-runtime --docker-opt=default-runtime=lynxi --docker-opt=registry-mirror=https://docker.996.ninja
./load_to_minikube.sh
minikube dashboard