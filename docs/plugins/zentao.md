# zentao Plugin

This plugin installs [zentao](https://zentao.net/) in an existing Kubernetes cluster by go client.

**Notes:**

- Zentao will be installed in K8S cluster, please prepare a k8s cluster before using zentao plugin.
  For local build, you can use `hack/e2e/e2e-up.sh` to create a k8s cluster via `Kind`.
- Currently, all fields list in the example config file below are required. You can modify them according to your needs.
- This plugin is not supported to run on `arm64` architecture now.

## Usage

```yaml
---
# core config
varFile: ''
toolFile: ''
state: # state config, backend can be local or s3
  backend: local
  options:
    stateFile: devstream.state

---
# plugins config
tools:
  # name of the tool
  - name: zentao
    # id of the tool instance
    instanceID: default
    # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool
    dependsOn: []
    # options for the plugin
    options:
      # namespace for zentao application
      namespace: 'zentao'
      # storageClassName used to match pv and pvc
      storageClassName: 'zentao-storage'
      # two PersistentVolumes for zentao and mysql should be specified
      persistentVolume:
        # name of zentao pv
        zentaoPVName: 'zentao-pv'
        # capacity of zentao pv
        zentaoPVCapacity: '1G'
        # name of mysql pv
        mysqlPVName: 'mysql-pv'
        # capacity of mysql pv
        mysqlPVCapacity: '1G'
      # two PersistentVolumeClaims for zentao and mysql should be specified
      persistentVolumeClaim:
        # name of zentao pvc
        zentaoPVCName: 'zentao-pvc'
        # capacity of zentao pvc
        zentaoPVCCapacity: '1G'
        # name of mysql pvc
        mysqlPVCName: 'mysql-pv'
        # capacity of mysql pvc
        mysqlPVCCapacity: '1G'
      # zentao application is deployed by K8S Deployment
      deployment:
        # name of zentao deployment
        name: 'zentao-dp'
        # number of application replica
        replicas: 3
        # zentao image
        image: 'easysoft/zentao:latest'
        # initial password name for mysql database, you can specify any name you like
        mysqlPasswdName: 'MYSQL_ROOT_PASSWORD'
        # initial password value for mysql database, you can specify any value you like
        mysqlPasswdValue: '1234567'
      # zentao application is exposed via K8S Service
      service:
        # name of zentao service
        name: 'zentao-svc'
        # nodePort of zentao service, currently zentao plugin only support `nodePort` type
        nodePort: 30081
```
