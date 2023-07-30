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
	"time"

	metricsv1 "github.com/crossplane-contrib/x-metrics/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

var i = 0
var generateNamespaceName = func() string {
	i = i + 1
	return fmt.Sprintf("test%d", i)
}

var crds = []apiextensions.CustomResourceDefinition{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nameas.testa.cloud",
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{

			Names: apiextensions.CustomResourceDefinitionNames{
				Categories: []string{
					"managed",
					"crda",
					"crdx",
				},
				Kind:   "NameA",
				Plural: "nameas",
			},
			Group: "testa.cloud",
			Scope: "Namespaced",
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Storage: true,
					Schema: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
						},
					},
				},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "namebs.testa.cloud",
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{

			Names: apiextensions.CustomResourceDefinitionNames{
				Categories: []string{
					"managed",
					"crdb",
				},
				Kind:   "NameB",
				Plural: "namebs",
			},
			Group: "testa.cloud",
			Scope: "Namespaced",
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name:    "v1beta1",
					Storage: true,
					Schema: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
						},
					},
				},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "namecs.testb.cloud",
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{

			Names: apiextensions.CustomResourceDefinitionNames{
				Categories: []string{
					"managed",
					"crdc",
					"crdx",
				},
				Kind:   "NameC",
				Plural: "namecs",
			},
			Group: "testb.cloud",
			Scope: "Namespaced",
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Storage: true,

					Schema: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
						},
					},
				},
				{
					Name: "v1beta1",
					Schema: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
						},
					},
				},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nameds.testb.cloud",
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{

			Names: apiextensions.CustomResourceDefinitionNames{
				Categories: []string{
					"managed",
					"crdd",
				},
				Kind:   "NameD",
				Plural: "nameds",
			},
			Group: "testb.cloud",
			Scope: "Namespaced",
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name:    "v2",
					Storage: true,
					Schema: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
						},
					},
				},
				{
					Name: "v1alpha1",
					Schema: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
						},
					},
				},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "namees.testc.cloud",
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{

			Names: apiextensions.CustomResourceDefinitionNames{
				Categories: []string{
					"managed",
					"crde",
				},
				Kind:   "NameE",
				Plural: "namees",
			},
			Group: "testc.cloud",
			Scope: "Namespaced",
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Storage: true,
					Schema: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
						},
					},
				},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "namefs.testc.cloud",
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{

			Names: apiextensions.CustomResourceDefinitionNames{
				Categories: []string{
					"managed",
					"crdf",
				},
				Kind:   "NameF",
				Plural: "namefs",
			},
			Group: "testc.cloud",
			Scope: "Cluster",
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Storage: true,
					Schema: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
						},
					},
				},
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "namegs.testc.cloud",
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{

			Names: apiextensions.CustomResourceDefinitionNames{
				Categories: []string{
					"managed",
					"crdg",
				},
				Kind:   "NameG",
				Plural: "namegs",
			},
			Group: "testc.cloud",
			Scope: "Namespaced",
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Storage: true,
					Schema: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Type: "object",
						},
					},
				},
			},
		},
	},
}

var _ = Describe("Managed Metrics", func() {
	Context("match", func() {
		It("Should install crds correctly", func() {
			ctx := context.Background()

			for _, crd := range crds {
				err := k8sClient.Create(ctx, &crd)
				Expect(err).NotTo(HaveOccurred(), "failed to create crds")
			}
		})
		It("Should select correct crds for single versions", func(ctx SpecContext) {
			mm.ResetRegister()
			metricsMemory = map[string]*MetricsMemory{}

			metricName := "cmetrica"
			metricNamespace := generateNamespaceName()

			mNamespace := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: metricNamespace,
				},
			}
			err2 := k8sClient.Create(ctx, &mNamespace)
			Expect(err2).NotTo(HaveOccurred(), "failed to create x-metrics namespace")
			matchName := "testa.cloud"
			metric := &metricsv1.Metric{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "metrics.crossplane.io/v1",
					Kind:       "ClusterMetric",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      metricName,
					Namespace: metricNamespace,
				},
				Spec: metricsv1.MetricSpec{
					MatchName: &matchName,
				},
			}

			Expect(k8sClient.Create(ctx, metric)).Should(Succeed())
			var mmMap map[string]schema.GroupVersionResource
			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).WithTimeout(time.Second * 20).Should(Equal(2))

			gvr1, ok := mmMap["testa_cloud_NameA_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameA v1")

			Expect(gvr1.Group).Should(Equal("testa.cloud"))
			Expect(gvr1.Version).Should(Equal("v1"))
			Expect(gvr1.Resource).Should(Equal("nameas"))

			gvr2, ok := mmMap["testa_cloud_NameB_v1beta1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameB v1beta1")

			Expect(gvr2.Group).Should(Equal("testa.cloud"))
			Expect(gvr2.Version).Should(Equal("v1beta1"))
			Expect(gvr2.Resource).Should(Equal("namebs"))

			Expect(k8sClient.Delete(ctx, metric)).Should(Succeed())
		}, SpecTimeout(time.Second*20))

		It("Should select correct crds for multiple versions", func(ctx SpecContext) {
			mm.ResetRegister()
			metricsMemory = map[string]*MetricsMemory{}

			metricName := "cmetrica"
			metricNamespace := generateNamespaceName()

			mNamespace := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: metricNamespace,
				},
			}
			err2 := k8sClient.Create(ctx, &mNamespace)
			Expect(err2).NotTo(HaveOccurred(), "failed to create x-metrics namespace")
			matchName := "testb.cloud"
			metric := &metricsv1.Metric{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "metrics.crossplane.io/v1",
					Kind:       "ClusterMetric",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      metricName,
					Namespace: metricNamespace,
				},
				Spec: metricsv1.MetricSpec{
					MatchName: &matchName,
				},
			}

			Expect(k8sClient.Create(ctx, metric)).Should(Succeed())
			var mmMap map[string]schema.GroupVersionResource
			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).WithTimeout(time.Second * 20).Should(Equal(4))

			gvr1, ok := mmMap["testb_cloud_NameC_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameC v1")

			Expect(gvr1.Group).Should(Equal("testb.cloud"))
			Expect(gvr1.Version).Should(Equal("v1"))
			Expect(gvr1.Resource).Should(Equal("namecs"))

			gvr2, ok := mmMap["testb_cloud_NameC_v1beta1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameC v1beta1")

			Expect(gvr2.Group).Should(Equal("testb.cloud"))
			Expect(gvr2.Version).Should(Equal("v1beta1"))
			Expect(gvr2.Resource).Should(Equal("namecs"))

			gvr3, ok := mmMap["testb_cloud_NameD_v2"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameD v2")

			Expect(gvr3.Group).Should(Equal("testb.cloud"))
			Expect(gvr3.Version).Should(Equal("v2"))
			Expect(gvr3.Resource).Should(Equal("nameds"))

			gvr4, ok := mmMap["testb_cloud_NameD_v1alpha1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameD v1alpha1")

			Expect(gvr4.Group).Should(Equal("testb.cloud"))
			Expect(gvr4.Version).Should(Equal("v1alpha1"))
			Expect(gvr4.Resource).Should(Equal("nameds"))

			Expect(k8sClient.Delete(ctx, metric)).Should(Succeed())
		}, SpecTimeout(time.Second*20))

		It("Should exclude crds", func(ctx SpecContext) {
			mm.ResetRegister()
			metricsMemory = map[string]*MetricsMemory{}

			metricName := "cmetrica"
			metricNamespace := generateNamespaceName()

			mNamespace := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: metricNamespace,
				},
			}
			err2 := k8sClient.Create(ctx, &mNamespace)
			Expect(err2).NotTo(HaveOccurred(), "failed to create x-metrics namespace")
			matchName := "testa.cloud"
			metric := &metricsv1.Metric{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "metrics.crossplane.io/v1",
					Kind:       "ClusterMetric",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      metricName,
					Namespace: metricNamespace,
				},
				Spec: metricsv1.MetricSpec{
					MatchName:    &matchName,
					ExcludeNames: &[]string{"namebs.testa.cloud"},
				},
			}

			Expect(k8sClient.Create(ctx, metric)).Should(Succeed())
			var mmMap map[string]schema.GroupVersionResource
			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).Should(Equal(1))

			gvr1, ok := mmMap["testa_cloud_NameA_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameA v1")

			Expect(gvr1.Group).Should(Equal("testa.cloud"))
			Expect(gvr1.Version).Should(Equal("v1"))
			Expect(gvr1.Resource).Should(Equal("nameas"))

			_, ok = mmMap["testa_cloud_NameB_v1beta1"]
			Expect(ok).Should(BeFalse(), "Should not have metrics for NameB v1beta1")

			Expect(k8sClient.Delete(ctx, metric)).Should(Succeed())
		}, SpecTimeout(time.Second*20))

		It("Should include crds", func(ctx SpecContext) {
			mm.ResetRegister()
			metricsMemory = map[string]*MetricsMemory{}

			metricName := "cmetrica"
			metricNamespace := generateNamespaceName()

			mNamespace := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: metricNamespace,
				},
			}
			err2 := k8sClient.Create(ctx, &mNamespace)
			Expect(err2).NotTo(HaveOccurred(), "failed to create x-metrics namespace")
			matchName := "testa.cloud"
			metric := &metricsv1.Metric{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "metrics.crossplane.io/v1",
					Kind:       "ClusterMetric",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      metricName,
					Namespace: metricNamespace,
				},
				Spec: metricsv1.MetricSpec{
					MatchName:    &matchName,
					IncludeNames: &[]string{"namecs.testb.cloud"},
				},
			}

			Expect(k8sClient.Create(ctx, metric)).Should(Succeed())
			var mmMap map[string]schema.GroupVersionResource
			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).Should(Equal(4))

			gvr1, ok := mmMap["testa_cloud_NameA_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameA v1")

			Expect(gvr1.Group).Should(Equal("testa.cloud"))
			Expect(gvr1.Version).Should(Equal("v1"))
			Expect(gvr1.Resource).Should(Equal("nameas"))

			gvr2, ok := mmMap["testa_cloud_NameB_v1beta1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameB v1beta1")

			Expect(gvr2.Group).Should(Equal("testa.cloud"))
			Expect(gvr2.Version).Should(Equal("v1beta1"))
			Expect(gvr2.Resource).Should(Equal("namebs"))

			gvr3, ok := mmMap["testb_cloud_NameC_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameC v1")

			Expect(gvr3.Group).Should(Equal("testb.cloud"))
			Expect(gvr3.Version).Should(Equal("v1"))
			Expect(gvr3.Resource).Should(Equal("namecs"))

			gvr4, ok := mmMap["testb_cloud_NameC_v1beta1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameC v1beta1")

			Expect(gvr4.Group).Should(Equal("testb.cloud"))
			Expect(gvr4.Version).Should(Equal("v1beta1"))
			Expect(gvr4.Resource).Should(Equal("namecs"))
			Expect(k8sClient.Delete(ctx, metric)).Should(Succeed())
		}, SpecTimeout(time.Second*20))

		It("Should match categories", func(ctx SpecContext) {
			mm.ResetRegister()
			metricsMemory = map[string]*MetricsMemory{}

			metricName := "cmetrica"
			metricNamespace := generateNamespaceName()

			mNamespace := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: metricNamespace,
				},
			}
			err2 := k8sClient.Create(ctx, &mNamespace)
			Expect(err2).NotTo(HaveOccurred(), "failed to create x-metrics namespace")

			metric := &metricsv1.Metric{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "metrics.crossplane.io/v1",
					Kind:       "ClusterMetric",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      metricName,
					Namespace: metricNamespace,
				},
				Spec: metricsv1.MetricSpec{
					Categories: &metricsv1.MetricCategory{
						Values: []string{
							"crda",
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, metric)).Should(Succeed())
			var mmMap map[string]schema.GroupVersionResource
			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).Should(Equal(1))

			gvr1, ok := mmMap["testa_cloud_NameA_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameA v1")

			Expect(gvr1.Group).Should(Equal("testa.cloud"))
			Expect(gvr1.Version).Should(Equal("v1"))
			Expect(gvr1.Resource).Should(Equal("nameas"))

			Expect(k8sClient.Delete(ctx, metric)).Should(Succeed())
		}, SpecTimeout(time.Second*20))

		It("Should match categories default to AND", func(ctx SpecContext) {
			mm.ResetRegister()
			metricsMemory = map[string]*MetricsMemory{}

			metricName := "cmetrica"
			metricNamespace := generateNamespaceName()

			mNamespace := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: metricNamespace,
				},
			}
			err2 := k8sClient.Create(ctx, &mNamespace)
			Expect(err2).NotTo(HaveOccurred(), "failed to create x-metrics namespace")

			metric := &metricsv1.Metric{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "metrics.crossplane.io/v1",
					Kind:       "ClusterMetric",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      metricName,
					Namespace: metricNamespace,
				},
				Spec: metricsv1.MetricSpec{
					Categories: &metricsv1.MetricCategory{
						Values: []string{
							"crda",
							"crdx",
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, metric)).Should(Succeed())
			var mmMap map[string]schema.GroupVersionResource
			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).Should(Equal(1))

			gvr1, ok := mmMap["testa_cloud_NameA_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameA v1")

			Expect(gvr1.Group).Should(Equal("testa.cloud"))
			Expect(gvr1.Version).Should(Equal("v1"))
			Expect(gvr1.Resource).Should(Equal("nameas"))

			Expect(k8sClient.Delete(ctx, metric)).Should(Succeed())
		}, SpecTimeout(time.Second*20))
		It("Should match categories OR", func(ctx SpecContext) {
			mm.ResetRegister()
			metricsMemory = map[string]*MetricsMemory{}

			metricName := "cmetrica"
			metricNamespace := generateNamespaceName()

			mNamespace := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: metricNamespace,
				},
			}
			err2 := k8sClient.Create(ctx, &mNamespace)
			Expect(err2).NotTo(HaveOccurred(), "failed to create x-metrics namespace")

			metric := &metricsv1.Metric{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "metrics.crossplane.io/v1",
					Kind:       "ClusterMetric",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      metricName,
					Namespace: metricNamespace,
				},
				Spec: metricsv1.MetricSpec{
					Categories: &metricsv1.MetricCategory{
						Values: []string{
							"crda",
							"crdx",
						},
						Join: "OR",
					},
				},
			}

			Expect(k8sClient.Create(ctx, metric)).Should(Succeed())
			var mmMap map[string]schema.GroupVersionResource
			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).Should(Equal(3))

			gvr1, ok := mmMap["testa_cloud_NameA_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameA v1")

			Expect(gvr1.Group).Should(Equal("testa.cloud"))
			Expect(gvr1.Version).Should(Equal("v1"))
			Expect(gvr1.Resource).Should(Equal("nameas"))

			gvr2, ok := mmMap["testb_cloud_NameC_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameC v1")

			Expect(gvr2.Group).Should(Equal("testb.cloud"))
			Expect(gvr2.Version).Should(Equal("v1"))
			Expect(gvr2.Resource).Should(Equal("namecs"))

			gvr3, ok := mmMap["testb_cloud_NameC_v1beta1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameC v1beta1")

			Expect(gvr3.Group).Should(Equal("testb.cloud"))
			Expect(gvr3.Version).Should(Equal("v1beta1"))
			Expect(gvr3.Resource).Should(Equal("namecs"))

			Expect(k8sClient.Delete(ctx, metric)).Should(Succeed())
		}, SpecTimeout(time.Second*20))

		It("Should prever matchNames", func(ctx SpecContext) {
			mm.ResetRegister()
			metricsMemory = map[string]*MetricsMemory{}

			metricName := "cmetrica"
			metricNamespace := generateNamespaceName()

			mNamespace := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: metricNamespace,
				},
			}
			err2 := k8sClient.Create(ctx, &mNamespace)
			Expect(err2).NotTo(HaveOccurred(), "failed to create x-metrics namespace")
			matchName := "testa.cloud"
			metric := &metricsv1.Metric{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "metrics.crossplane.io/v1",
					Kind:       "ClusterMetric",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      metricName,
					Namespace: metricNamespace,
				},
				Spec: metricsv1.MetricSpec{
					MatchName: &matchName,
					Categories: &metricsv1.MetricCategory{
						Values: []string{
							"crdb",
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, metric)).Should(Succeed())
			var mmMap map[string]schema.GroupVersionResource
			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).Should(Equal(2))

			gvr1, ok := mmMap["testa_cloud_NameA_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameA v1")

			Expect(gvr1.Group).Should(Equal("testa.cloud"))
			Expect(gvr1.Version).Should(Equal("v1"))
			Expect(gvr1.Resource).Should(Equal("nameas"))

			gvr2, ok := mmMap["testa_cloud_NameB_v1beta1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameB v1beta1")

			Expect(gvr2.Group).Should(Equal("testa.cloud"))
			Expect(gvr2.Version).Should(Equal("v1beta1"))
			Expect(gvr2.Resource).Should(Equal("namebs"))

			Expect(k8sClient.Delete(ctx, metric)).Should(Succeed())
		}, SpecTimeout(time.Second*20))

		It("Should add finalizer", func(ctx SpecContext) {
			mm.ResetRegister()
			metricsMemory = map[string]*MetricsMemory{}

			metricName := "cmetrica"
			metricNamespace := generateNamespaceName()

			mNamespace := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: metricNamespace,
				},
			}
			err2 := k8sClient.Create(ctx, &mNamespace)
			Expect(err2).NotTo(HaveOccurred(), "failed to create x-metrics namespace")
			matchName := "testa.cloud"
			metric := &metricsv1.Metric{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "metrics.crossplane.io/v1",
					Kind:       "ClusterMetric",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      metricName,
					Namespace: metricNamespace,
				},
				Spec: metricsv1.MetricSpec{
					MatchName: &matchName,
					Categories: &metricsv1.MetricCategory{
						Values: []string{
							"crdb",
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, metric)).Should(Succeed())

			Eventually(func() int {
				var met metricsv1.Metric
				metRequest := types.NamespacedName{
					Name:      metricName,
					Namespace: metricNamespace,
				}
				err := k8sClient.Get(ctx, metRequest, &met)
				if err != nil {
					return -1
				}
				fin := met.Finalizers
				return len(fin)

			}).Should(Equal(1))

			Expect(k8sClient.Delete(ctx, metric)).Should(Succeed())
		}, SpecTimeout(time.Second*20))

		It("Should remove metric on update", func(ctx SpecContext) {
			mm.ResetRegister()
			metricsMemory = map[string]*MetricsMemory{}

			metricName := "cmetrica"
			metricNamespace := generateNamespaceName()

			mNamespace := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: metricNamespace,
				},
			}
			err2 := k8sClient.Create(ctx, &mNamespace)
			Expect(err2).NotTo(HaveOccurred(), "failed to create x-metrics namespace")

			metric := &metricsv1.Metric{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "metrics.crossplane.io/v1",
					Kind:       "ClusterMetric",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      metricName,
					Namespace: metricNamespace,
				},
				Spec: metricsv1.MetricSpec{
					Categories: &metricsv1.MetricCategory{
						Values: []string{
							"crda",
							"crdx",
						},
						Join: "OR",
					},
				},
			}

			Expect(k8sClient.Create(ctx, metric)).Should(Succeed())
			var mmMap map[string]schema.GroupVersionResource
			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).Should(Equal(3))

			gvr1, ok := mmMap["testa_cloud_NameA_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameA v1")

			Expect(gvr1.Group).Should(Equal("testa.cloud"))
			Expect(gvr1.Version).Should(Equal("v1"))
			Expect(gvr1.Resource).Should(Equal("nameas"))

			gvr2, ok := mmMap["testb_cloud_NameC_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameC v1")

			Expect(gvr2.Group).Should(Equal("testb.cloud"))
			Expect(gvr2.Version).Should(Equal("v1"))
			Expect(gvr2.Resource).Should(Equal("namecs"))

			gvr3, ok := mmMap["testb_cloud_NameC_v1beta1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameC v1beta1")

			Expect(gvr3.Group).Should(Equal("testb.cloud"))
			Expect(gvr3.Version).Should(Equal("v1beta1"))
			Expect(gvr3.Resource).Should(Equal("namecs"))

			var met metricsv1.Metric
			metRequest := types.NamespacedName{
				Name:      metricName,
				Namespace: metricNamespace,
			}
			Expect(k8sClient.Get(ctx, metRequest, &met)).Should(Succeed())
			met.Spec.Categories.Join = "AND"
			Expect(k8sClient.Update(ctx, &met)).Should(Succeed())

			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).Should(Equal(1))

			gvr4, ok := mmMap["testa_cloud_NameA_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameA v1")

			Expect(gvr4.Group).Should(Equal("testa.cloud"))
			Expect(gvr4.Version).Should(Equal("v1"))
			Expect(gvr4.Resource).Should(Equal("nameas"))

			Expect(k8sClient.Delete(ctx, metric)).Should(Succeed())
		}, SpecTimeout(time.Second*20))

		It("Should remove all metric on update", func(ctx SpecContext) {
			mm.ResetRegister()
			metricsMemory = map[string]*MetricsMemory{}

			metricName := "cmetrica"
			metricNamespace := generateNamespaceName()

			mNamespace := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: metricNamespace,
				},
			}
			err2 := k8sClient.Create(ctx, &mNamespace)
			Expect(err2).NotTo(HaveOccurred(), "failed to create x-metrics namespace")

			metric := &metricsv1.Metric{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "metrics.crossplane.io/v1",
					Kind:       "ClusterMetric",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      metricName,
					Namespace: metricNamespace,
				},
				Spec: metricsv1.MetricSpec{
					Categories: &metricsv1.MetricCategory{
						Values: []string{
							"crdy",
							"crdx",
						},
						Join: "OR",
					},
				},
			}

			Expect(k8sClient.Create(ctx, metric)).Should(Succeed())
			var mmMap map[string]schema.GroupVersionResource
			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).Should(Equal(3))

			gvr1, ok := mmMap["testa_cloud_NameA_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameA v1")

			Expect(gvr1.Group).Should(Equal("testa.cloud"))
			Expect(gvr1.Version).Should(Equal("v1"))
			Expect(gvr1.Resource).Should(Equal("nameas"))

			gvr2, ok := mmMap["testb_cloud_NameC_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameC v1")

			Expect(gvr2.Group).Should(Equal("testb.cloud"))
			Expect(gvr2.Version).Should(Equal("v1"))
			Expect(gvr2.Resource).Should(Equal("namecs"))

			gvr3, ok := mmMap["testb_cloud_NameC_v1beta1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameC v1beta1")

			Expect(gvr3.Group).Should(Equal("testb.cloud"))
			Expect(gvr3.Version).Should(Equal("v1beta1"))
			Expect(gvr3.Resource).Should(Equal("namecs"))

			var met metricsv1.Metric
			metRequest := types.NamespacedName{
				Name:      metricName,
				Namespace: metricNamespace,
			}
			Expect(k8sClient.Get(ctx, metRequest, &met)).Should(Succeed())
			met.Spec.Categories.Join = "AND"
			Expect(k8sClient.Update(ctx, &met)).Should(Succeed())

			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).Should(Equal(0))

			Expect(k8sClient.Delete(ctx, metric)).Should(Succeed())
		}, SpecTimeout(time.Second*20))

		It("Should add only namespaces resources", func(ctx SpecContext) {
			mm.ResetRegister()
			metricsMemory = map[string]*MetricsMemory{}

			metricName := "cmetrica"
			metricNamespace := generateNamespaceName()

			mNamespace := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: metricNamespace,
				},
			}
			matchName := "testc.cloud"
			err2 := k8sClient.Create(ctx, &mNamespace)
			Expect(err2).NotTo(HaveOccurred(), "failed to create x-metrics namespace")

			metric := &metricsv1.Metric{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "metrics.crossplane.io/v1",
					Kind:       "ClusterMetric",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      metricName,
					Namespace: metricNamespace,
				},
				Spec: metricsv1.MetricSpec{
					MatchName: &matchName,
				},
			}

			Expect(k8sClient.Create(ctx, metric)).Should(Succeed())
			var mmMap map[string]schema.GroupVersionResource
			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).Should(Equal(2))

			gvr1, ok := mmMap["testc_cloud_NameE_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameE v1")

			Expect(gvr1.Group).Should(Equal("testc.cloud"))
			Expect(gvr1.Version).Should(Equal("v1"))
			Expect(gvr1.Resource).Should(Equal("namees"))

			gvr2, ok := mmMap["testc_cloud_NameG_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameG v1")

			Expect(gvr2.Group).Should(Equal("testc.cloud"))
			Expect(gvr2.Version).Should(Equal("v1"))
			Expect(gvr2.Resource).Should(Equal("namegs"))

			Expect(k8sClient.Delete(ctx, metric)).Should(Succeed())
		}, SpecTimeout(time.Second*20))

		It("Should add metrics only once", func(ctx SpecContext) {
			mm.ResetRegister()
			metricsMemory = map[string]*MetricsMemory{}

			metricName := "cmetrica"
			metric2Name := "cmetricb"
			metricNamespace := generateNamespaceName()

			mNamespace := corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: metricNamespace,
				},
			}
			err2 := k8sClient.Create(ctx, &mNamespace)
			Expect(err2).NotTo(HaveOccurred(), "failed to create x-metrics namespace")

			metric := &metricsv1.Metric{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "metrics.crossplane.io/v1",
					Kind:       "ClusterMetric",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      metricName,
					Namespace: metricNamespace,
				},
				Spec: metricsv1.MetricSpec{
					Categories: &metricsv1.MetricCategory{
						Values: []string{
							"crda",
							"crdx",
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, metric)).Should(Succeed())
			var mmMap map[string]schema.GroupVersionResource
			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).Should(Equal(1))

			gvr1, ok := mmMap["testa_cloud_NameA_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameA v1")

			Expect(gvr1.Group).Should(Equal("testa.cloud"))
			Expect(gvr1.Version).Should(Equal("v1"))
			Expect(gvr1.Resource).Should(Equal("nameas"))

			metric2 := &metricsv1.Metric{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "metrics.crossplane.io/v1",
					Kind:       "ClusterMetric",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      metric2Name,
					Namespace: metricNamespace,
				},
				Spec: metricsv1.MetricSpec{
					Categories: &metricsv1.MetricCategory{
						Values: []string{
							"crda",
							"crdx",
						},
						Join: "OR",
					},
				},
			}

			Expect(k8sClient.Create(ctx, metric2)).Should(Succeed())

			// var mmMap map[string]schema.GroupVersionResource
			Eventually(func() int {
				mmMap = mm.GetRegister()
				return len(mmMap)
			}).Should(Equal(3))

			gvr2, ok := mmMap["testa_cloud_NameA_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameA v1")

			Expect(gvr2.Group).Should(Equal("testa.cloud"))
			Expect(gvr2.Version).Should(Equal("v1"))
			Expect(gvr2.Resource).Should(Equal("nameas"))

			gvr3, ok := mmMap["testb_cloud_NameC_v1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameC v1")

			Expect(gvr3.Group).Should(Equal("testb.cloud"))
			Expect(gvr3.Version).Should(Equal("v1"))
			Expect(gvr3.Resource).Should(Equal("namecs"))

			gvr4, ok := mmMap["testb_cloud_NameC_v1beta1"]
			Expect(ok).Should(BeTrue(), "Should have metrics for NameC v1beta1")

			Expect(gvr4.Group).Should(Equal("testb.cloud"))
			Expect(gvr4.Version).Should(Equal("v1beta1"))
			Expect(gvr4.Resource).Should(Equal("namecs"))

			numCalls := mm.GetNumOfCalls()

			for _, v := range numCalls {
				Expect(v).Should(Equal(1))
			}

			Expect(k8sClient.Delete(ctx, metric)).Should(Succeed())
			Expect(k8sClient.Delete(ctx, metric2)).Should(Succeed())
		}, SpecTimeout(time.Second*20))

		It("Should delete crds correctly", func() {
			ctx := context.Background()

			for _, crd := range crds {
				err := k8sClient.Delete(ctx, &crd)
				Expect(err).NotTo(HaveOccurred(), "failed to delete crds")
			}
		})
	})
})
