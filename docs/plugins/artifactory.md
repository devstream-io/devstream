# artifactory Plugin

This plugin installs [artifactory](https://jfrog.com/artifactory/) in an existing Kubernetes cluster using the Helm chart.

## Usage

### Test/Local Dev Environment

If you want to **test the plugin locally**， The following `valuesYaml` configuration can be used

```yaml
valuesYaml: |
  artifactory:
    service:
      type: NodePort
    nodePort: 30002
  nginx:
    enabled: false
```

In this configuration

- Postgresql dependencies are automatically created.
- local disks on machines in the cluster are defaulted used for data mounting.
- Using `nodePort` to expose service, You can access `artifactory` by domain `http://{{k8s node IP}}:30002`. The default account name and password are admin/password (please replace the default account password in the production environment).

### Production Environment

#### External Storage

- PostgreSQL: Set the `database.url` to Postgresql's address. More info can be found in [Config](https://www.jfrog.com/confluence/display/JFROG/Configuring+the+Database).

#### Disk Storage

You can set `customVolumes` and `customVolumeMounts` for this service. More info can be found in [Config](https://www.jfrog.com/confluence/display/JFROG/Configuring+the+Filestore).

#### Network Config

This plugin support`Ingress`, `ClusterIP`, `NodePort` and `LoadBalancer` , You can give choice to your needs.

### Config

The following content is an example of the "tool file".

For more information on the main config, the tool file and the var file of DevStream, see [Core Concepts Overview](../core-concepts/core-concepts.md#1-config) and [DevStream Configuration](../core-concepts/config.md).

```yaml
--8<-- "artifactory.yaml"
```

#### Default Configs

| key                | default value           | description                                        |
| ----               | ----                    | ----                                               |
| chart.chartPath    | ""                      | local chart path                                   |
| chart.chartName    | jfrog/artifactory       | chart name                                         |
| chart.timeout      | 10m                     | this config will wait 10 minutes to deploy         |
| chart.releaseName  | artifactory             | helm release name                                  |
| chart.upgradeCRDs  | true                    | default update CRD config                          |
| chart.wait         | true                    | whether to wait until installation is complete     |
| chart.namespace    | artifactory             | namespace where helm to deploy                     |
| repo.url           | https://charts.jfrog.io | offical helm repo address                          |
| repo.name          | jfrog                   | helm repo name                                     |
Currently, except for `valuesYaml` and default configs, all the parameters in the example above are mandatory.
