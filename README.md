# x-metrics
`x-metrics` generates Prometheus metrics for a range of Crossplane resources, encompassing Providers, Compositions, Claims, Managed Resources, etc. These metrics provide comprehensive insights, including details such as the last transition time, creation timestamp, readiness status, and more. Access to these metrics is available via an exposed endpoint.

based on [Crossplane Intro and Deep Dive - the Cloud Native Control Plane Framework](https://youtu.be/5WRkVUlEgHI?t=1793)

## Prerequisites

[Helm](https://helm.sh) must be installed to use the charts.
Please refer to Helm's [documentation](https://helm.sh/docs/) to get started.

## Installation

Install the Helm chart:
```console
helm repo add x-metrics https://crossplane-contrib.github.io/x-metrics
helm install x-metrics x-metrics/x-metrics --namespace x-metrics --create-namespace --wait
```

## Usage

To access the metrics with the default setting trough the endpoint:
1. Port-forward the services
```console
kubectl -n x-metrics port-forward svc/x-metrics 8080:8080
```
2. In your browser navigate to: http://127.0.0.1:8080/x-metrics

3. To generate metrics, apply one of the CRDs under the `examples/` folder:
```console
kubectl appy -f examples/iam-metric.yaml
```
4. Refresh the browser to see the metrics populate.

## Licensing

| Property                       | Function              | Repository  |
|--------------------------------|-----------------------|-------------|
| metrics                        | metrics               | [xp-state-metrics](https://github.com/chlunde/xp-state-metrics) |
| managed-metrics                | metrics               | [managed-metrics](https://github.com/dkb-bank/managed-metrics) |
| kube-state-metrics             | metrics               | [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics)
