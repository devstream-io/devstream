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

## Deployment

### Step1: Prepare a Kubernetes Cluster
- If you already have a kubernetes cluster, ignore this step. 
- If not, you can use `hack/e2e/e2e-up.sh` to create a k8s cluster via `Kind` as test environment.
  ```shell
  bash hack/e2e/e2e-up.sh
  ```

### Step2: Create Zentao Application via Config File
- Create a zentao config file following the usage example above.
  ```shell
  ./dtm apply -f zentao.yaml --debug
  ```

### Step3: Initialize Zentao Application
- Visit `http://NodeIP:NodePort`("NodeIP" and "NodePort" are Kubernets node IP and node port) to start the initialization process. Press `Start Installation` button to the next step.
![](zentao/zentao-welcome.jpg)

- You don't need to do anything about the system check and it's done automatically. If there are system check items that do not pass, please make sure that the previous operation is correct. If it still doesn't work, create an issue to track your problem.
![](zentao/zentao-systemCheck.jpg)

- Fill in database password filed with `options.deployment.mysqlPasswdValue` which was set previously in `zentao.yaml`.
![](zentao/zentao-configuration.jpg)

- If everything proceeds successfully, you will see the Zendo introduction.
![](zentao/zentao-intro.jpg)

- Fill in your company name and create an administrator account.
![](zentao/zentao-account.jpg)

- Now, the Zendo application has been successfully deployed.
![](zentao/zentao-web.jpg)
