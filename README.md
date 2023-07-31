# x-metrics

based on [Crossplane Intro and Deep Dive - the Cloud Native Control Plane Framework](https://youtu.be/5WRkVUlEgHI?t=1793)


## Usage

[Helm](https://helm.sh) must be installed to use the charts.
Please refer to Helm's [documentation](https://helm.sh/docs/) to get started.

Once Helm is set up properly, add the repo as follows:

```console
helm repo add x-metrics https://crossplane-contrib.github.io/x-metrics
helm install x-metrics x-metrics/x-metrics --namespace x-metrics --create-namespace --wait
```

 ## Licensing

| Property                       | Function              | Repository  |
|--------------------------------|-----------------------|-------------|
| metrics                        | metrics               | [xp-state-metrics](https://github.com/chlunde/xp-state-metrics) |
| managed-metrics                | metrics               | [managed-metrics](https://github.com/dkb-bank/managed-metrics) |
| kube-state-metrics             | metrics               | [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics)
