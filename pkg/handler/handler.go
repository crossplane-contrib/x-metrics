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

package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kube-state-metrics/v2/pkg/metric"
	"sigs.k8s.io/controller-runtime/pkg/log"

	store "github.com/crossplane-contrib/x-metrics/pkg/store"
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IManagedMetricsHandler interface {
	ServeHTTP(writer http.ResponseWriter, r *http.Request)
	RegisterAndAddMetricStoreForGVR(ctx context.Context, metricName string, gvr schema.GroupVersionResource, namespace string) chan struct{}
	RemoveMetricStore(name string)
}

type ManagedMetricsHandler struct {
	metricsWriter   map[string]store.IXMetricsStore
	Client          dynamic.Interface
	callbacks       map[string]func() (schema.GroupVersionResource, int)
	newStoreHandler func([]string, func(interface{}) []metric.FamilyInterface, context.Context, dynamic.Interface, string, schema.GroupVersionResource, string) store.IXMetricsStore
}

type InfoMappings struct {
	FieldPath string
	Label     string
}
type crossplaneStatus struct {
	ready      float64
	synced     float64
	readyTime  time.Time
	syncedTime time.Time
}

func NewManagedMetricsHandler(dc dynamic.Interface) ManagedMetricsHandler {
	return ManagedMetricsHandler{
		metricsWriter:   map[string]store.IXMetricsStore{},
		Client:          dc,
		callbacks:       map[string]func() (schema.GroupVersionResource, int){},
		newStoreHandler: store.NewXMetricsStore,
	}
}

func NewManagedMetricsHandlerWithStore(dc dynamic.Interface, storeHandler func([]string, func(interface{}) []metric.FamilyInterface, context.Context, dynamic.Interface, string, schema.GroupVersionResource, string) store.IXMetricsStore) ManagedMetricsHandler {
	return ManagedMetricsHandler{
		metricsWriter:   map[string]store.IXMetricsStore{},
		Client:          dc,
		callbacks:       map[string]func() (schema.GroupVersionResource, int){},
		newStoreHandler: storeHandler,
	}
}

// nolint: errcheck
func (m *ManagedMetricsHandler) ServeHTTP(writer http.ResponseWriter, r *http.Request) {

	totalCount := 0
	for _, c := range m.callbacks {
		_, count := c()
		totalCount += count
	}
	for _, w := range m.metricsWriter {
		w.WriteAll(writer)
	}

	writer.Write([]byte("# TYPE x_metric_resources_count_total gauge\n# HELP x_metric_resources_count_total A metric to count all resources\n"))
	writer.Write([]byte("x_metric_resources_count_total "))
	writer.Write([]byte(strconv.Itoa(totalCount)))
	writer.Write([]byte{'\n'})

	if closer, ok := writer.(io.Closer); ok {
		closer.Close()
	}
}

func (m *ManagedMetricsHandler) RegisterAndAddMetricStoreForGVR(ctx context.Context, metricName string, gvr schema.GroupVersionResource, namespace string) chan struct{} {
	reflectorStore, channel := m.registerMetricStoreForGVR(ctx, metricName, gvr, namespace)
	m.addMetricStore(metricName, reflectorStore)
	return channel
}

func (m *ManagedMetricsHandler) addMetricStore(name string, metricStore store.IXMetricsStore) {
	m.metricsWriter[name] = metricStore
}

func (m *ManagedMetricsHandler) RemoveMetricStore(name string) {
	metricsStore := m.metricsWriter[name]
	callbackUid := metricsStore.GetCallbacUid()
	delete(m.callbacks, callbackUid)
	delete(m.metricsWriter, name)
}

func (m *ManagedMetricsHandler) registerMetricStoreForGVR(ctx context.Context, metricName string, gvr schema.GroupVersionResource, namespace string) (store.IXMetricsStore, chan struct{}) {

	log := log.FromContext(ctx)

	if namespace != "" {
		metricName = GetValidLabel(namespace + "_" + metricName)
	}
	headers := []string{
		"# TYPE %s gauge\n# HELP %s A metrics series for each object",
		"# TYPE %s_created gauge\n# HELP %s_created Unix creation timestamp",
		"# TYPE %s_labels gauge\n# HELP %s_labels Labels from the kubernetes object",
		"# TYPE %s_info gauge\n# HELP %s_info A metrics series exposing parameters as labels",
		"# TYPE %s_ready gauge\n# HELP %s_ready A metrics series mapping the Ready status condition to a value (True=1,False=0,other=-1)",
		"# TYPE %s_ready_time gauge\n# HELP %s_ready_time Unix timestamp of last ready change",
		"# TYPE %s_synced gauge\n# HELP %s_synced A metrics series mapping the Synced status condition to a value (True=1,False=0,other=-1)",
		"# TYPE %s_synced_time gauge\n# HELP %s_synced_time Unix timestamp of last synced change",
	}
	for i, hfmt := range headers {
		headers[i] = fmt.Sprintf(hfmt, metricName, metricName)
	}
	labelKeys := []string{"name"}
	labelValues := func(obj *unstructured.Unstructured) []string {
		return []string{obj.GetName()}
	}

	if namespace != "" {
		labelKeys = append(labelKeys, "namespace")
		labelValues = func(obj *unstructured.Unstructured) []string {
			return []string{obj.GetName(), obj.GetNamespace()}
		}
	}
	reflectorStore := m.newStoreHandler(headers, func(objAny any) []metric.FamilyInterface {
		obj := objAny.(*unstructured.Unstructured)
		paved := fieldpath.Pave(obj.Object)
		o := metric.Family{
			Name: metricName,
			Metrics: []*metric.Metric{
				{
					LabelKeys:   labelKeys,
					LabelValues: labelValues(obj),
					Value:       1,
				},
			},
		}

		families := []metric.FamilyInterface{&o}

		created := metric.Family{
			Name: metricName + "_created",
			Metrics: []*metric.Metric{
				{
					LabelKeys:   labelKeys,
					LabelValues: labelValues(obj),
					Value:       float64(obj.GetCreationTimestamp().Unix()),
				},
			},
		}
		families = append(families, &created)

		labels := metric.Family{
			Name: metricName + "_labels",
			Metrics: []*metric.Metric{
				{
					LabelKeys:   labelKeys,
					LabelValues: labelValues(obj),
					Value:       1,
				},
			},
		}
		for k, v := range obj.GetLabels() {
			labels.Metrics[0].LabelKeys = append(labels.Metrics[0].LabelKeys, "label_"+GetValidLabel(k))
			labels.Metrics[0].LabelValues = append(labels.Metrics[0].LabelValues, v)
		}
		families = append(families, &labels)

		mappings := []InfoMappings{}

		var infoKeys, infoValues []string
		for _, m := range mappings {
			val, _ := paved.GetString(m.FieldPath)
			infoKeys = append(infoKeys, m.Label)
			infoValues = append(infoValues, val)
		}

		o_info := metric.Family{
			Name: metricName + "_info",
			Metrics: []*metric.Metric{
				{
					LabelKeys:   append(labelKeys, infoKeys...),
					LabelValues: append(labelValues(obj), infoValues...),
					Value:       1,
				},
			},
		}

		families = append(families, &o_info)

		status := getCrossplaneStatus(obj)
		o_ready := metric.Family{
			Name: metricName + "_ready",
			Metrics: []*metric.Metric{
				{
					LabelKeys:   labelKeys,
					LabelValues: labelValues(obj),
					Value:       status.ready,
				},
			},
		}

		families = append(families, o_ready)

		o_ready_time := metric.Family{
			Name: metricName + "_ready_time",
			Metrics: []*metric.Metric{
				{
					LabelKeys:   labelKeys,
					LabelValues: labelValues(obj),
					Value:       float64(status.readyTime.Unix()),
				},
			},
		}

		families = append(families, o_ready_time)

		o_synced := metric.Family{
			Name: metricName + "_synced",
			Metrics: []*metric.Metric{
				{
					LabelKeys:   labelKeys,
					LabelValues: labelValues(obj),
					Value:       status.synced,
				},
			},
		}

		families = append(families, o_synced)

		o_synced_time := metric.Family{
			Name: metricName + "_synced_time",
			Metrics: []*metric.Metric{
				{
					LabelKeys:   labelKeys,
					LabelValues: labelValues(obj),
					Value:       float64(status.syncedTime.Unix()),
				},
			},
		}

		families = append(families, o_synced_time)

		return families
	}, ctx, m.Client, namespace, gvr, metricName)
	callbackUid, countCallBack := reflectorStore.GetCallback()
	m.callbacks[callbackUid] = countCallBack
	lw := cache.ListWatch{
		ListFunc: func(opt metav1.ListOptions) (runtime.Object, error) {
			o, err := m.Client.Resource(gvr).Namespace(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				log.Info("err listing")
			}
			return o, err
		},
		WatchFunc: func(ops metav1.ListOptions) (watch.Interface, error) {
			return m.Client.Resource(gvr).Namespace(namespace).Watch(ctx, ops)
		},
	}

	re := cache.NewReflector(&lw, &unstructured.Unstructured{}, reflectorStore, 0)

	channel := make(chan struct{})
	go re.Run(channel)

	return reflectorStore, channel
}

func GetValidLabel(name string) string {

	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return r
		case r >= 'a' && r <= 'z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '-',
			r == '_',
			r == '.',
			r == '/':
			return '_'
		}
		return -1
	}, name)
}

func statusToPrometheusValue(s xpv1.ConditionedStatus, typ xpv1.ConditionType) float64 {
	switch s.GetCondition(typ).Status {
	case "True":
		return 1
	case "False":
		return 0
	default:
		return -1
	}
}

func getCrossplaneStatus(u *unstructured.Unstructured) crossplaneStatus {
	conditioned := xpv1.ConditionedStatus{}
	_ = fieldpath.Pave(u.Object).GetValueInto("status", &conditioned)

	return crossplaneStatus{
		ready:      statusToPrometheusValue(conditioned, xpv1.TypeReady),
		synced:     statusToPrometheusValue(conditioned, xpv1.TypeSynced),
		readyTime:  conditioned.GetCondition(xpv1.TypeReady).LastTransitionTime.Time,
		syncedTime: conditioned.GetCondition(xpv1.TypeSynced).LastTransitionTime.Time,
	}
}
