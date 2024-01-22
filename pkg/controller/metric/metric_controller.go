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

package controllers

import (
	"context"
	"fmt"
	"regexp"
	"time"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	metricsv1 "github.com/crossplane-contrib/x-metrics/api/v1"
	xmetrics "github.com/crossplane-contrib/x-metrics/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// MetricReconciler reconciles a Metric object
type MetricReconciler struct {
	client.Client
	Kind      string
	Scheme    *runtime.Scheme
	MmHandler xmetrics.IManagedMetricsHandler
}

type Resource struct {
	Group      string  `json:"group"`
	Version    string  `json:"version"`
	Kind       string  `json:"kind"`
	Resource   string  `json:"resource"`
	MetricName string  `json:"metricName"`
	Namespace  *string `json:"namespace,omitempty"`
}

type MetricsDefinition struct {
	Channels map[string]*CloseChannel
}

type CloseChannel struct {
	Channel chan struct{}
	Closed  bool
}

type MetricsMemory struct {
	Channel  *CloseChannel
	Consumer map[string]struct{}
}

const (
	finalizerName = "metrics.crossplane.io/finalizer"
)

var (
	metricsMemory = map[string]*MetricsMemory{}
)

func (r *MetricReconciler) newReconciler() (client.Object, error) {
	metricGVK := metricsv1.GroupVersion.WithKind(r.Kind)
	ro, err := r.Scheme.New(metricGVK)
	if err != nil {
		return nil, err
	}
	return ro.(client.Object), nil
}

// +kubebuilder:rbac:groups=metrics.crossplane.io,resources=metrics,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metrics.crossplane.io,resources=metrics/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metrics.crossplane.io,resources=metrics/finalizers,verbs=update
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list
func (r *MetricReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx)
	metric, err := r.newReconciler()
	if err != nil {
		log.Error(err, "Unrecognised metric type")
		return ctrl.Result{}, nil
	}
	if err := r.Get(ctx, req.NamespacedName, metric); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var namespaced bool
	switch metric.(type) {
	case *metricsv1.Metric:
		namespaced = true
	case *metricsv1.ClusterMetric:
		namespaced = false
	default:
		log.Error(fmt.Errorf("unexpected metric type: %t", metric), "Not retrying.")
		return ctrl.Result{}, nil
	}
	var currentConsumerName string
	var currentNamespace string
	if namespaced {
		currentConsumerName = metric.GetNamespace() + "::" + metric.GetName()
		currentNamespace = metric.GetNamespace()

	} else {
		currentConsumerName = metric.GetName()
	}
	currentMetrics := getCurrentMetrics(currentConsumerName)

	objectMeta, metricSpec, metricStatus, _ := getSpecAndStatus(metric)

	// Add a finalizer, if the object is not already marked for deletion
	if objectMeta.DeletionTimestamp.IsZero() {

		if !controllerutil.ContainsFinalizer(metric, finalizerName) {
			controllerutil.AddFinalizer(metric, finalizerName)
			if err := r.Update(ctx, metric); err != nil {
				return ctrl.Result{}, nil
			}

		}
	} else {
		// If the object is marked for deletion, run the cleanup, if a finaliser is set
		if controllerutil.ContainsFinalizer(metric, finalizerName) {
			cleanupMetrics(r.MmHandler, currentMetrics, currentConsumerName)
			controllerutil.RemoveFinalizer(metric, finalizerName)
			if err := r.Update(ctx, metric); err != nil {
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, nil
		}

	}

	resourceList, err := r.getGVRForMetric(ctx, metricSpec, namespaced)

	addR, _, deleteR := r.getResources(ctx, &currentMetrics, resourceList)

	var statusMetrics []metricsv1.WatchedResource
	if metricStatus.WatchedResources != nil {
		statusMetrics = *metricStatus.WatchedResources
	} else {
		statusMetrics = []metricsv1.WatchedResource{}
	}
	if len(addR) > 0 {
		for _, v := range addR {
			metricName := v.MetricName
			if _, ok := metricsMemory[metricName]; !ok {
				gvr := schema.GroupVersionResource{
					Group:    v.Group,
					Version:  v.Version,
					Resource: v.Resource,
				}
				channel := r.MmHandler.RegisterAndAddMetricStoreForGVR(ctx, metricName, gvr, currentNamespace)
				metricsMemory[metricName] = &MetricsMemory{
					Consumer: map[string]struct{}{
						currentConsumerName: {},
					},
					Channel: &CloseChannel{
						Channel: channel,
						Closed:  false,
					},
				}
			} else {
				metricsMemory[metricName].Consumer[currentConsumerName] = struct{}{}
			}

			statusMetrics = append(statusMetrics, metricsv1.WatchedResource{
				Kind:       v.Kind,
				Group:      v.Group,
				Version:    v.Version,
				MetricName: &metricName,
			})

		}
	}

	if len(deleteR) > 0 {
		cleanupMetrics(r.MmHandler, deleteR, currentConsumerName)
		statusMetrics = filterDeletedMetrics(&statusMetrics, &deleteR)
	}
	if err != nil {
		log.Error(err, "unable to get resources")
	}
	metricStatus.WatchedResources = &statusMetrics
	if err := r.Client.Status().Update(ctx, metric); err != nil {
		log.Error(err, "unable to update metric status")
	}

	duration := time.Minute * 5

	return ctrl.Result{RequeueAfter: duration}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *MetricReconciler) SetupWithManager(mgr ctrl.Manager) error {
	reconcilerType, err := r.newReconciler()
	if err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(reconcilerType).
		Complete(r)
}

func (r *MetricReconciler) getGVRForMetric(ctx context.Context, metric *metricsv1.MetricSpec, namespaced bool) (*map[string]Resource, error) {

	list := map[string]Resource{}

	var crds apiextensions.CustomResourceDefinitionList
	options := client.ListOptions{}
	if err := r.Client.List(ctx, &crds, &options); err != nil {
		return nil, err
	}
	if metric.MatchName != nil || metric.Categories != nil {
		for _, crd := range crds.Items {
			name := crd.GetName()
			match := true
			if metric.MatchName != nil {
				match, _ = regexp.MatchString(*metric.MatchName, name)
			} else if metric.Categories != nil {
				crdCategories := crd.Spec.Names.Categories
				match = matchesCategories(crdCategories, metric.Categories.Values, metric.Categories.Join)
			}
			inNameList := inList(metric.IncludeNames, name)
			inExcludeList := inList(metric.ExcludeNames, name)
			inNamespace := isNamespaced(&crd)
			// if we need a gvr for a metrics resource, we only watch namespaced resources
			if (match || inNameList) && !inExcludeList && (namespaced == inNamespace || !namespaced) {
				for _, version := range crd.Spec.Versions {
					if version.Storage {
						metricName := xmetrics.GetValidLabel(crd.Spec.Group + "_" + crd.Spec.Names.Kind + "_" + version.Name)
						list[metricName] = Resource{
							Group:      crd.Spec.Group,
							Version:    version.Name,
							Resource:   crd.Spec.Names.Plural,
							Kind:       crd.Spec.Names.Kind,
							MetricName: metricName,
						}
					}
				}
			}
		}
	}
	return &list, nil
}

func matchesCategories(current []string, wanted []string, joinType metricsv1.MetricJoin) bool {
	contains := false
	for _, w := range wanted {
		for _, c := range current {
			if w == c {
				contains = true
			}
			// if the jointype is "or", one single match is sufficient to return true
			if joinType == metricsv1.JoinOr && contains {
				return true
			}
		}
		// if the jointype is "AND", one missed wanted string is sufficient to return false
		if joinType == metricsv1.JoinAnd && !contains {
			return false
		}
	}
	return contains
}

func getCurrentMetrics(cMetricName string) []string {

	currentMetrics := []string{}
	for metricName, metric := range metricsMemory {
		for clusertMetricName := range metric.Consumer {
			if clusertMetricName == cMetricName {
				currentMetrics = append(currentMetrics, metricName)
			}
		}
	}
	return currentMetrics
}

func cleanupMetrics(handler xmetrics.IManagedMetricsHandler, metrics []string, currentConsumer string) {
	for _, metricName := range metrics {
		if metric, ok := metricsMemory[metricName]; ok {
			if !metric.Channel.Closed {
				delete(metric.Consumer, currentConsumer)
				if len(metric.Consumer) == 0 {
					close(metric.Channel.Channel)
					metric.Channel.Closed = true
					handler.RemoveMetricStore(metricName)
					delete(metricsMemory, metricName)
				}
			}

		}
	}
}

func (r *MetricReconciler) getResources(ctx context.Context, currentMetics *[]string, resouces *map[string]Resource) (map[string]Resource, []string, []string) {
	add := map[string]Resource{}
	current := []string{}
	delete := []string{}

	currentMap := map[string]struct{}{}

	if currentMetics != nil {
		for _, key := range *currentMetics {

			currentMap[key] = struct{}{}
			if _, ok := (*resouces)[key]; !ok {
				delete = append(delete, key)
			}
		}
	}

	for k, resource := range *resouces {
		_, ok := currentMap[k]

		if ok {
			current = append(current, k)
		} else {
			add[k] = resource
		}
	}
	return add, current, delete
}

func filterDeletedMetrics(metrics *[]metricsv1.WatchedResource, delteList *[]string) []metricsv1.WatchedResource {

	newMetricList := []metricsv1.WatchedResource{}

mainloop:
	for _, v := range *metrics {
		for _, id := range *delteList {
			if id == *v.MetricName {
				continue mainloop
			}
		}
		newMetricList = append(newMetricList, v)

	}
	return newMetricList
}

func inList(list *[]string, value string) bool {
	if list != nil {
		for _, v := range *list {
			if value == v {
				return true
			}
		}
	}
	return false
}

func isNamespaced(crd *apiextensions.CustomResourceDefinition) bool {
	return crd.Spec.Scope == "Namespaced"
}

func getSpecAndStatus(metric client.Object) (*metav1.ObjectMeta, *metricsv1.MetricSpec, *metricsv1.MetricStatus, error) {
	switch t := metric.(type) {
	case *metricsv1.Metric:
		return &t.ObjectMeta, &t.Spec, &t.Status, nil
	case *metricsv1.ClusterMetric:
		return &t.ObjectMeta, &t.Spec, &t.Status, nil
	default:
		return nil, nil, nil, fmt.Errorf("not an metric type: %t", t)
	}
}
