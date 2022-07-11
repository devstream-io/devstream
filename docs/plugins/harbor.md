# harbor Plugin

This plugin installs [harbor](https://goharbor.io/) in an existing Kubernetes cluster using the Helm chart.

## Usage

##### External storage

- Postgresql: Set the `database.type` to `external` and fill the information in `database.external` section.
- Redis: Set the `redis.type` to `external` and fill the information in `redis.external` section.

##### Disk storage

Please be sure to change the `storageClass` in the options of the config to an existing StorageClass. You can refer to this [document](https://github.com/goharbor/harbor-helm#configure-how-to-persist-data)

#### Network config

This plugin support `Ingress`, `ClusterIP`, `NodePort`, `LoadBalancer` ，You can give choice to your needs.

##### TLS config
- Use self-signed certificates
  1. Set `tls.enabled` to `true` and edit the corresponding domain name `externalURL`
  2. Copy the self-signed certificate stored in the Pod `harbor-core` directory `/etc/core/ca` to your own PC
  3. Trust the certificate on your own host
- Using public Certificates
  1. Add the public certificate as Secret
  2. Set `tls.enabled` to `true` and edit the corresponding domain name `externalURL`
  3. Configure `tls.secretName` to use the public certificate


### Test/Local Dev Environment

If you want **test plugin locally**，you can just use default params

- Postgresql and Redis dependencies are automatically created
- By default, local disks on machines in the cluster are used for data mounting
- Set the value of `values_yaml` in the configuration as follows. Use `nodePort` to provide the harbor service externally
```yaml
        values_yaml: |
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
- Now you can access harbor through the domain name `http://{{k8s node IP}}:30002`. The default account name and password are admin/Harbor12345 (please replace the default account password in production environment)

### Config

```yaml
tools:
# name of the tool
- name: harbor
  # id of the tool instance
  instanceID: default
  # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ ]
  # options for the plugin
  options:
    create_namespace: true
    repo:
      name: harbor
      # url of the Helm repo, use self host helm config beacuse official helm does'nt support namespace config
      url: https://helm.goharbor.io
    # Helm chart information
    chart:
      # name of the chart
      chart_name: harbor/harbor
      # k8s namespace where Tekton will be installed
      namespace: harbor
      # release name of the chart
      release_name: harbor
      # whether to wait for the release to be deployed or not
      wait: true
      # the time to wait for any individual Kubernetes operation (like Jobs for hooks). This defaults to 5m0s
      timeout: 10m
      # whether to perform a CRD upgrade during installation
      upgradeCRDs: true
      values_yaml: |
          trivy:
            enabled: false
```

Currently, all the parameters in the example above are mandatory.

