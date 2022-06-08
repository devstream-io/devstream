# zentao 插件

该插件将通过`client-go`库在 Kubernetes 集群中安装[禅道应用](https://zentao.net/).

**注意:**

- 使用该插件之前，用户需要准备好 Kubernetes 集群，并正确设置 KUBECONFIG 配置。本地测试时，用户可以通过`hack/e2e/e2e-up.sh`创建临时 Kubernetes 集群。
- 用户可以根据自己的需求修改以下配置文件中的字段，所有列出的字段都需要被设置。
- 该插件暂时不支持运行在`arm64`系统架构；`amd64`架构已经过测试运行正常。

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
