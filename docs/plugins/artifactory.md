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

Currently, except for `values_yaml`, all the parameters in the example above are mandatory.
