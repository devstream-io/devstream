# `jenkins` Plugin

This plugin installs [jenkins](https://jenkins.io) in an existing Kubernetes cluster using the Helm chart.

## Usage

NOTICE:

Create the data directory manually:

```bash
mkdir -p ~/data/jenkins-volumes/
chown -R 1000:1000 ~/data/jenkins-volumes/
```

```yaml
tools:
# name of the instance with jenkins
- name: jenkins-dev
  # name of the plugin
  plugin: jenkins
  # optional; if specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ "TOOL1_NAME.TOOL1_PLUGIN", "TOOL2_NAME.TOOL2_PLUGIN" ]
  # options for the plugin
  options:
    # need to create the namespace or not, default: false
    create_namespace: false
    # Helm repo information
    repo:
      # name of the Helm repo
      name: jenkins
      # url of the Helm repo
      url: https://charts.jenkins.io
    # Helm chart information
    chart:
      # name of the chart
      chart_name: jenkins/jenkins
      # release name of the chart
      release_name: dev
      # k8s namespace where jenkins will be installed
      namespace: jenkins
      # whether to wait for the release to be deployed or not
      wait: true
      # the time to wait for any individual Kubernetes operation (like Jobs for hooks). This defaults to 5m0s
      timeout: 5m
      # whether to perform a CRD upgrade during installation
      upgradeCRDs: true
      # custom configuration (Optional). You can refer to [Jenkins values.yaml](https://github.com/jenkinsci/helm-charts/blob/main/charts/jenkins/values.yaml)
      values_yaml: |
        persistence:
          storageClass: jenkins-pv
        serviceAccount:
          create: false
          name: jenkins

```

Currently, all the parameters in the example above are mandatory.
