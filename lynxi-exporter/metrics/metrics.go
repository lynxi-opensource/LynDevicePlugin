// Package metrics 定义和记录Prometheus指标
package metrics

// Recorder 描述了一个定义和记录Prometheus指标的对象需要提供的方法
type Recorder interface {
	Record() error
}
