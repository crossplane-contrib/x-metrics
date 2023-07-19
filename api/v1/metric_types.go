/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MetricSpec defines the desired state of Metric
type MetricSpec struct {

	// MatchName is a string to match CRDs with names that match this string
	MatchName *string `json:"matchName,omitempty"`
	// IncludeNames lists crds that should be added to metrics
	IncludeNames *[]string `json:"includeNames,omitempty"`
	// ExcludeNames lists crds that should not be added to metrics. If they are added by other metrics objects, they are not excluded explicitly
	ExcludeNames *[]string `json:"excludeNames,omitempty"`

	// Categories contains an object to add metrics for crds by crd category. Categories are only evaluated, if MatchName is nil
	Categories *MetricCategory `json:"categories,omitempty"`
}

// MetricStatus defines the observed state of Metric
type MetricStatus struct {
	MetricBaseName   *string            `json:"metricBaseName,omitempty"`
	WatchedResources *[]WatchedResource `json:"watchedResources,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Namespaced

// Metric is the Schema for the Metrics API
type Metric struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MetricSpec   `json:"spec,omitempty"`
	Status MetricStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MetricList contains a list of Metric
type MetricList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Metric `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Metric{}, &MetricList{})
}

// +kubebuilder:validation:Enum=AND;OR
type MetricJoin string

const (
	JoinAnd MetricJoin = "AND"
	JoinOr  MetricJoin = "OR"
)

type MetricCategory struct {
	// Values is a list of strings
	Values []string `json:"values"`

	// Join decides if a single value in Values is needed or all to add metrics for the corresponding resource
	// +kubebuilder:default:=AND
	Join MetricJoin `json:"join,omitempty"`
}

type WatchedResource struct {
	Group      string  `json:"group"`
	Version    string  `json:"version"`
	Kind       string  `json:"kind"`
	Namespace  *string `json:"namespace,omitempty"`
	MetricName *string `json:"metricName,omitempty"`
}
