package handler_test

import (
	"context"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/crossplane-contrib/x-metrics/pkg/handler"
	store_test "github.com/crossplane-contrib/x-metrics/pkg/handler/mock"
)

var _ = Describe("Handler", func() {
	Context("servehttp", func() {
		It("Should write metric for total number of objects", func() {
			testEnv := &envtest.Environment{
				CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "package", "crds")},
				ErrorIfCRDPathMissing: true,
			}

			cfg, _ := testEnv.Start()

			dc, _ := dynamic.NewForConfig(cfg)
			handler := handler.NewManagedMetricsHandlerWithStore(dc, store_test.NewXMetricsStoreMockGenerator(0, ""))
			w := store_test.ResponseWriterMock{}
			handler.ServeHTTP(&w, nil)

			Expect(w.Data).Should(ContainSubstring("# TYPE x_metric_resources_count_total gauge"))
			Expect(w.Data).Should(ContainSubstring("# HELP x_metric_resources_count_total A metric to count all resources"))
			Expect(w.Data).Should(ContainSubstring("x_metric_resources_count_total 0"))
		})
		It("Should register a metric", func() {
			testEnv := &envtest.Environment{
				CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "package", "crds")},
				ErrorIfCRDPathMissing: true,
			}

			cfg, _ := testEnv.Start()
			ctx, cancel := context.WithCancel(context.TODO())
			defer func() {
				cancel()
			}()
			dc, _ := dynamic.NewForConfig(cfg)
			handler := handler.NewManagedMetricsHandlerWithStore(dc, store_test.NewXMetricsStoreMockGenerator(5, "Test"))
			handler.RegisterAndAddMetricStoreForGVR(ctx, "test", schema.GroupVersionResource{
				Group:    "test",
				Version:  "v1",
				Resource: "object",
			}, "")
			w := store_test.ResponseWriterMock{}
			handler.ServeHTTP(&w, nil)

			Expect(w.Data).Should(ContainSubstring("# TYPE x_metric_resources_count_total gauge"))
			Expect(w.Data).Should(ContainSubstring("# HELP x_metric_resources_count_total A metric to count all resources"))
			Expect(w.Data).Should(ContainSubstring("x_metric_resources_count_total 5"))
		})
	})
})
