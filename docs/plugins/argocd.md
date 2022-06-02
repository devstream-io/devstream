# argocd Plugin

This plugin installs [ArgoCD](https://argoproj.github.io/cd/) in an existing Kubernetes cluster using the Helm chart.

## Usage

```yaml
--8<-- "argocd.yaml"
```

Currently, except for `values_yaml`, all the parameters in the example above are mandatory.
