package store_test

import (
	"context"
	"io"
	"net/http"

	"github.com/crossplane-contrib/x-metrics/pkg/store"
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/kube-state-metrics/v2/pkg/metric"
)

type XMetricsStoreMock struct {
	Num       int
	WriteData string
	uid       string
}

// nolint: errcheck
func (x *XMetricsStoreMock) WriteAll(w io.Writer) {
	w.Write([]byte(x.WriteData))
}

func (x *XMetricsStoreMock) GetCallback() (string, func() (schema.GroupVersionResource, int)) {

	uid := uuid.New().String()
	x.uid = uid
	return uid, func() (schema.GroupVersionResource, int) {
		return schema.GroupVersionResource{}, x.Num
	}
}

func (x *XMetricsStoreMock) GetCallbacUid() string {
	return ""
}
func (x *XMetricsStoreMock) Add(interface{}) error {
	return nil
}
func (s *XMetricsStoreMock) Delete(obj interface{}) error {
	return nil
}

func (s *XMetricsStoreMock) Update(obj interface{}) error {
	// TODO: For now, just call Add, in the future one could check if the resource version changed?
	return nil
}

// List implements the List method of the store interface.
func (s *XMetricsStoreMock) List() []interface{} {

	return nil
}

// ListKeys implements the ListKeys method of the store interface.
func (s *XMetricsStoreMock) ListKeys() []string {
	return nil
}

// Get implements the Get method of the store interface.
func (s *XMetricsStoreMock) Get(obj interface{}) (item interface{}, exists bool, err error) {
	return nil, false, nil
}

// GetByKey implements the GetByKey method of the store interface.
func (s *XMetricsStoreMock) GetByKey(key string) (item interface{}, exists bool, err error) {
	return nil, false, nil
}

// Replace will delete the contents of the store, using instead the
// given list.
func (s *XMetricsStoreMock) Replace(list []interface{}, _ string) error {
	return nil
}

// Resync implements the Resync method of the store interface.
func (s *XMetricsStoreMock) Resync() error {
	return nil
}

func NewXMetricsStoreMockGenerator(num int, data string) func([]string, func(interface{}) []metric.FamilyInterface, context.Context, dynamic.Interface, string, schema.GroupVersionResource, string) store.IXMetricsStore {

	return func(headers []string, generateFunc func(interface{}) []metric.FamilyInterface, ctx context.Context, client dynamic.Interface, namespace string, gvr schema.GroupVersionResource, metricName string) store.IXMetricsStore {

		store := &XMetricsStoreMock{
			Num:       num,
			WriteData: data,
		}

		return store

	}
}

type ResponseWriterMock struct {
	Data string
}

func (w *ResponseWriterMock) Header() http.Header {
	return http.Header{}
}

func (w *ResponseWriterMock) Write(d []byte) (int, error) {

	w.Data += string(d)
	return 0, nil
}

func (w *ResponseWriterMock) WriteHeader(statusCode int) {

}
