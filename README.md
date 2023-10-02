# x-metrics
`x-metrics` generates Prometheus metrics for a range of Crossplane resources, encompassing Providers, Compositions, Claims, Managed Resources, etc. These metrics provide comprehensive insights, including details such as the last transition time, creation timestamp, readiness status, and more. Access to these metrics is available via an exposed endpoint.

based on [Crossplane Intro and Deep Dive - the Cloud Native Control Plane Framework](https://youtu.be/5WRkVUlEgHI?t=1793)


## Usage

[Helm](https://helm.sh) must be installed to use the charts.
Please refer to Helm's [documentation](https://helm.sh/docs/) to get started.

Once Helm is set up properly, add the repo as follows:

```console
helm repo add x-metrics https://crossplane-contrib.github.io/x-metrics
helm install x-metrics x-metrics/x-metrics --namespace x-metrics --create-namespace --wait
```

Access the metrics trough the endpoint. For the default service settings:
```cosole
kubect -n x-metrics port-forward svc/x-metrics 8080:8080
```
In your browser navigate to: http://127.0.0.1:8080/x-metrics

Metrics should be empty. To generate metrics, apply one of the CRDs under the `examples/` folder:
```console
kubectl appy -f examples/iam-metric.yaml
```
Refresh the browser to see the metrics populate.

## Licensing

| Property                       | Function              | Repository  |
|--------------------------------|-----------------------|-------------|
| metrics                        | metrics               | [xp-state-metrics](https://github.com/chlunde/xp-state-metrics) |
| managed-metrics                | metrics               | [managed-metrics](https://github.com/dkb-bank/managed-metrics) |
| kube-state-metrics             | metrics               | [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics)
