# zentao(禅道)插件

该插件将通过`client-go`库在 Kubernetes 集群中安装[禅道应用](https://zentao.net/).

**注意:**

- 使用该插件之前，用户需要准备好 Kubernetes 集群，并正确设置 KUBECONFIG 配置。本地测试时，用户可以通过`hack/e2e/e2e-up.sh`创建临时 Kubernetes 集群。
- 用户可以根据自己的需求修改以下配置文件中的字段，所有列出的字段都需要被设置。
- 该插件暂时不支持运行在`arm64` CPU 架构上；在`amd64` CPU 架构上运行正常。

## 用法示例

下面的配置文件展示的是"tool file"的内容。

关于更多关于DevStream的主配置、tool file、var file的信息，请阅读[核心概念概览](../core-concepts/overview.zh.md)和[DevStream配置](../core-concepts/config.zh.md).

```yaml
config:
  state:
    backend: local
    options:
      stateFile: devstream.state

tools:
  # name of the tool
  - name: zentao
    # id of the tool instance
    instanceID: default
    # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool
    dependsOn: []
    # options for the plugin
    options:
      # namespace for ZenTao application
      namespace: 'zentao'
      # storageClassName used to match pv and pvc
      storageClassName: 'zentao-storage'
      # two PersistentVolumes for ZenTao and mysql should be specified
      persistentVolumes:
          # name of ZenTao pv
        - pvName: "zentao-pv"
          # capacity of ZenTao pv
          pvCapacity: "1G"
          # name of mysql pv
        - pvName: "mysql-pv"
          # capacity of mysql pv
          pvCapacity: "1G"
      # two PersistentVolumeClaims for ZenTao and mysql should be specified
      persistentVolumeClaims:
          # name of ZenTao pvc
        - pvcName: "zentao-pvc"
          # capacity of ZenTao pvc
          pvcCapacity: "1G"
          # name of mysql pvc
        - pvcName: "mysql-pvc"
          # capacity of mysql pvc
          pvcCapacity: "1G"
      # ZenTao application is deployed by K8s Deployment
      deployment:
        # name of ZenTao deployment
        name: 'zentao-dp'
        # number of application replica
        replicas: 1
        # ZenTao image
        image: 'easysoft/zentao:latest'
        envs:
          - key: 'MYSQL_ROOT_PASSWORD'
            # initial password value for mysql database, you can specify any value you like
            value: '123456'
      # ZenTao application is exposed via K8s Service
      service:
        # name of ZenTao service
        name: 'zentao-svc'
        # nodePort of ZenTao service, currently ZenTao plugin only support `nodePort` type
        nodePort: 30081
```

## 部署

### 第一步：准备一个 Kubernetes 集群

- 如果已经部署 Kubernetes 集群环境，可以忽略第一步。 
- 如果尚未部署 Kubernetes 集群环境，可以使用脚本`hack/e2e/e2e-up.sh`来创建一个 Kubernetes 测试环境，该脚本会基于`Kind`创建集群。
  
```shell
bash hack/e2e/e2e-up.sh
```

### 第二步：利用配置文件创建禅道应用

- 根据上面提供的示例用法创建一个`zentao.yaml`.

```shell
./dtm apply -f zentao.yaml --debug
```

### 第三步：初始化禅道应用

1. 浏览器访问`http://NodeIP:NodePort`("NodeIP"与"NodePort"对应 Kubernets 节点IP和节点端口)来初始化禅道服务。点击`开始安装`按钮进入下一步。
![](zentao/zentao-welcome.jpg)

2. 针对系统检查步骤，用户不需要做任何操作，系统检查会自动完成。如果存在未通过的检查项，请确保之前的部署流程均按照上述文档指示进行。如果仍不能解决问题，请在社区新建 issue 来跟踪相关问题。
![](zentao/zentao-systemCheck.jpg)

3. 数据库部署步骤只需在`数据库密码`栏填入之前设置的`options.deployment.mysqlPasswdValue`值。
![](zentao/zentao-configuration.jpg)

4. 如果上述所有操作均成功完成，这时将展示禅道应用介绍界面。
![](zentao/zentao-intro.jpg)

5. 填写用户的公司信息并创建管理员账户，然后点击`保存`。
![](zentao/zentao-account.jpg)

6. 至此，禅道服务就部署成功啦！
![](zentao/zentao-web.jpg)
