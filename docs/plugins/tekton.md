# tekton Plugin
This plugin installs [tekton]("https://tekton.dev/") in an existing Kubernetes cluster using the Helm chart.

## Usage

```yaml
tools:
  # name of the tool
  - name: tekton
    # id of the tool instance
    instanceID: default
    # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
    dependsOn: [ ]
    # options for the plugin
    options:
      # need to create the namespace or not, default: false
      create_namespace: true
      repo:
        # name of the Helm repo
        name: tekton
        # url of the Helm repo, use self host helm config beacuse official helm does'nt support namespace config
        url: https://steinliber.github.io/tekton-helm-chart/
      # Helm chart information
      chart:
        # name of the chart
        chart_name: tekton/tekton-pipeline
        # k8s namespace where Tekton will be installed
        namespace: tekton
        # release name of the chart
        release_name: tekton
        # whether to wait for the release to be deployed or not
        wait: true
        # the time to wait for any individual Kubernetes operation (like Jobs for hooks). This defaults to 5m0s
        timeout: 5m
        # whether to perform a CRD upgrade during installation
        upgradeCRDs: true
        values_yaml: |
          serviceaccount:
            enabled: true
```

Currently, except for `values_yaml`, all the parameters in the example above are mandatory.
