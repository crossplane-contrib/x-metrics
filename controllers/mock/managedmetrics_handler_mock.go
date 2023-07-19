package xmetrics

import (
	"context"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ManagedMetricsHandlerMock struct {
	register      map[string]schema.GroupVersionResource
	multipleCalls map[string]int
}

func NewManagedMetricsHandlerMock() ManagedMetricsHandlerMock {
	return ManagedMetricsHandlerMock{
		register:      map[string]schema.GroupVersionResource{},
		multipleCalls: map[string]int{},
	}
}

func (m *ManagedMetricsHandlerMock) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	for k := range m.register {
		writer.Write([]byte(k + ";"))
	}
}

func (m *ManagedMetricsHandlerMock) RegisterAndAddMetricStoreForGVR(ctx context.Context, metricName string, gvr schema.GroupVersionResource, namespace string) chan struct{} {
	if _, ok := m.register[metricName]; ok {
		m.multipleCalls[metricName] = m.multipleCalls[metricName] + 1
	} else {
		m.multipleCalls[metricName] = 1
	}
	m.register[metricName] = gvr
	return make(chan struct{})
}

func (m *ManagedMetricsHandlerMock) GetRegister() map[string]schema.GroupVersionResource {
	return m.register
}

func (m *ManagedMetricsHandlerMock) GetNumOfCalls() map[string]int {
	return m.multipleCalls
}

func (m *ManagedMetricsHandlerMock) ResetRegister() {
	m.register = map[string]schema.GroupVersionResource{}
	m.multipleCalls = map[string]int{}
}
func (m *ManagedMetricsHandlerMock) RemoveMetricStore(name string) {
	delete(m.register, name)
}
