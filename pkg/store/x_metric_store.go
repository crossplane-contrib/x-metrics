package store

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kube-state-metrics/v2/pkg/metric"
	metricsstore "k8s.io/kube-state-metrics/v2/pkg/metrics_store"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type IXMetricsStore interface {
	cache.Store
	WriteAll(io.Writer)
	GetCallbacUid() string
	GetCallback() (string, func() (schema.GroupVersionResource, int))
}

type XMetricsStore struct {
	metricStore metricsstore.MetricsStore
	counter     int
	gvr         schema.GroupVersionResource
	mutex       sync.RWMutex
	metricaName string
	callbackUid string
}

func NewXMetricsStore(headers []string, generateFunc func(interface{}) []metric.FamilyInterface, ctx context.Context, client dynamic.Interface, namespace string, gvr schema.GroupVersionResource, metricName string) IXMetricsStore {

	store := &XMetricsStore{
		metricStore: *metricsstore.NewMetricsStore(headers, generateFunc),
		gvr:         gvr,
		counter:     0,
		metricaName: metricName,
	}

	store.init(ctx, client, namespace, gvr)
	return store

}

func (s *XMetricsStore) init(ctx context.Context, client dynamic.Interface, namespace string, gvr schema.GroupVersionResource) {
	log := log.FromContext(ctx)
	o, err := client.Resource(gvr).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Error(err, "err listing")
	} else {
		s.counter = len(o.Items)
	}
}

func (s *XMetricsStore) Add(obj interface{}) error {
	s.counter++
	return s.metricStore.Add(obj)
}

func (s *XMetricsStore) Update(obj interface{}) error {
	// TODO: For now, just call Add, in the future one could check if the resource version changed?
	return s.metricStore.Update(obj)
}

func (s *XMetricsStore) Delete(obj interface{}) error {
	s.counter--
	return s.metricStore.Delete(obj)
}

// List implements the List method of the store interface.
func (s *XMetricsStore) List() []interface{} {

	return s.metricStore.List()
}

// ListKeys implements the ListKeys method of the store interface.
func (s *XMetricsStore) ListKeys() []string {
	return s.metricStore.ListKeys()
}

// Get implements the Get method of the store interface.
func (s *XMetricsStore) Get(obj interface{}) (item interface{}, exists bool, err error) {
	return s.metricStore.Get(obj)
}

// GetByKey implements the GetByKey method of the store interface.
func (s *XMetricsStore) GetByKey(key string) (item interface{}, exists bool, err error) {
	return s.metricStore.GetByKey(key)
}

// Replace will delete the contents of the store, using instead the
// given list.
func (s *XMetricsStore) Replace(list []interface{}, _ string) error {
	return s.metricStore.Replace(list, "")
}

// Resync implements the Resync method of the store interface.
func (s *XMetricsStore) Resync() error {
	return s.metricStore.Resync()
}

// WriteAll writes all metrics of the store into the given writer, zipped with the
// help text of each metric family.
func (s *XMetricsStore) WriteAll(w io.Writer) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.metricStore.WriteAll(w)
	s.writeCount(w)
}

// nolint: errcheck
func (s *XMetricsStore) writeCount(w io.Writer) {
	metricName := fmt.Sprintf("%s_resource_count", s.metricaName)
	w.Write([]byte(fmt.Sprintf("# TYPE %[1]s gauge\n# HELP %[1]s A metrics series objects to count objects of %[2]s\n", metricName, s.metricaName)))
	w.Write([]byte(metricName))
	w.Write([]byte(" "))
	w.Write([]byte(strconv.Itoa(s.counter)))
	w.Write([]byte{'\n'})
}

func (s *XMetricsStore) ConterAndType() (schema.GroupVersionResource, int) {
	return s.gvr, s.counter
}

func (s *XMetricsStore) GetCallback() (string, func() (schema.GroupVersionResource, int)) {
	uid := uuid.New().String()
	s.callbackUid = uid
	return uid, func() (schema.GroupVersionResource, int) {
		return s.ConterAndType()
	}
}

func (s *XMetricsStore) GetCallbacUid() string {
	return s.callbackUid
}
