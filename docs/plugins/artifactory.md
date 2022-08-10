# artifactory Plugin

This plugin installs [artifactory](https://jfrog.com/artifactory/) in an existing Kubernetes cluster using the Helm chart.

## Usage

### Test/Local Dev Environment

If you want to **test the plugin locally**， The following `values_yaml` configuration can be used

```yaml
values_yaml: |
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

```yaml
--8<-- "artifactory.yaml"
```

#### Default Configs

| key                | default value           | description                                        |
| ----               | ----                    | ----                                               |
| chart.chart_name   | jfrog/artifactory       | chart name                                         |
| chart.timeout      | 10m                     | this config will wait 10 minutes to deploy         |
| chart.release_name | artifactory             | helm release name                                  |
| chart.upgradeCRDs  | true                    | default update CRD config                          |
| chart.wait         | true                    | whether to wait until installation is complete     |
| repo.url           | https://charts.jfrog.io | offical helm repo address                          |
| repo.name          | jfrog                   | helm repo name                                     |
| create_namespace   | false                   | whether to create namespace if namespace not eixst |

Currently, except for `values_yaml` and default configs, all the parameters in the example above are mandatory.
