mkdir release -p
go build -o release/lynxi-k8s-device-plugin lyndeviceplugin/lynxi-k8s-device-plugin
go build -o release/lynxi-exporter lyndeviceplugin/lynxi-exporter