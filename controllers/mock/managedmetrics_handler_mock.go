/*
Copyright 2023 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
		_, err := writer.Write([]byte(k + ";"))
		if err != nil {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}
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
