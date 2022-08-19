# harbor Plugin

This plugin installs [harbor](https://goharbor.io/) in an existing Kubernetes cluster using the Helm chart.

## Usage

### Test/Local Dev Environment

If you want to **test the plugin locally**， The following `values_yaml` configuration can be used

```yaml
values_yaml: |
  externalURL: http://127.0.0.1
  expose:
    type: nodePort
    tls:
      enabled: false
  chartmuseum:
    enabled: false
  clair:
    enabled: false
  notary:
    enabled: false
  trivy:
    enabled: false
```

In this configuration

- Postgresql and Redis dependencies are automatically created.
- local disks on machines in the cluster are defaulted used for data mounting.
- Only the `harbor` main program is installed, not the rest of the plug-ins.
- Using `nodePort` to expose service, You can access `harbor` by domain `http://{{k8s node IP}}:30002`. The default account name and password are admin/Harbor12345 (please replace the default account password in the production environment).

### Production Environment

Most of Harbor's components are stateless. So we can simply increase the replica of the pods to make sure the components are distributed to multiple worker nodes and leverage the `Service` mechanism of `Kubernetes` to ensure connectivity across pods.

#### External Storage

> Both Postgresql and Redis have to be installed for the harbor.

- Postgresql: Set the `database.type` to `external` and fill the information in `database.external` section.
- Redis: Set the `redis.type` to `external` and fill the information in `redis.external` section.

#### Disk Storage

Please be sure to change the `storageClass` in the options of the config to an existing StorageClass, You can refer to this [document](https://github.com/goharbor/harbor-helm#configure-how-to-persist-data).

#### Network Config

This plugin support `Ingress`, `ClusterIP`, `NodePort`, `LoadBalancer` , You can give choice to your needs.

#### TLS Config

- Use self-signed certificates
  1. Set `tls.enabled` to `true` and edit the corresponding domain name `externalURL`.
  2. Copy the self-signed certificate stored in the Pod `harbor-core` directory `/etc/core/ca` to your own PC.
  3. Trust the certificate on your host.
- Using public Certificates
  1. Add the public certificate as Secret.
  2. Set `tls.enabled` to `true` and edit the corresponding domain name `externalURL`.
  3. Configure `tls.secretName` to use the public certificate.

### Config

```yaml
--8<-- "harbor.yaml"
```

#### Default Configs

| key                | default value            | description                                        |
| ----               | ----                     | ----                                               |
| chart.chart_name   | harbor/harbor            | chart name                                         |
| chart.timeout      | 10m                      | this config will wait 10 minutes to deploy         |
| chart.release_name | harbor                   | helm release name                                  |
| chart.upgradeCRDs  | true                     | default update CRD config                          |
| chart.wait         | true                     | whether to wait until installation is complete     |
| chart.namespace    | harbor                   | namespace where helm to deploy                     |
| repo.url           | https://helm.goharbor.io | helm official repo address                         |
| repo.name          | harbor                   | helm repo name                                     |
| create_namespace   | true                     | make sure namespace exist                          |

Currently, except for `values_yaml` and default configs, all the parameters in the example above are mandatory.
