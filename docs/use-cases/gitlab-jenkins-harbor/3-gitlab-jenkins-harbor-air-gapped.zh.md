# 离线环境快速部署 GitLab + Jenkins + Harbor 工具链

在“[这个文档](./2-gitlab-jenkins-harbor.zh.md)”里我们介绍了怎样通过 DevStream 在本地快速部署 `GitLab + Jenkins + Harbor` 工具链。

但是如果你的服务器是离线的，你只有一台可以访问互联网的 PC，这台 PC 可以通过企业内部网络访问到你要用来部署 `GitLab + Jenkins + Harbor` 工具链服务器，类似下图这样：

<figure markdown>
  ![GitLab token](./gitlab-jenkins-harbor-air-gapped/air-gapped-env.png){ width="500" }
  <figcaption></figcaption>
</figure>

这时候，你就需要用到 DevStream 的工具链离线部署能力了。

## 1、下载 dtm 和 DevStream Plugins

首先你需要下载 DevStream 的命令行（CLI）工具 `dtm` 和所需的 DevStream 插件（plugins）。

### 1.1、下载 dtm

你可以参考[这个文档](../../install.zh.md)下载 dtm。

唯一需要注意的是，下载完之后，请记得将 dtm 传输到你需要使用它的机器上。

### 1.2、下载 plugins

继续在你的 PC 上执行如下命令来下载 DevStream plugins：

```shell
dtm init --download-only --plugins="gitlab-ce-docker, helm-installer" -d=plugins
```

这条命令执行成功后，你可以在本地 plugins 目录下看到如下文件：

```shell
$ ls plugins/
gitlab-ce-docker-linux-amd64_0.10.3.md5
helm-installer-linux-amd64_0.10.3.md5
gitlab-ce-docker-linux-amd64_0.10.3.so
helm-installer-linux-amd64_0.10.3.so
```

## 2、下载镜像

因为 DevStream 需要使用容器化方式部署 GitLab、Jenkins 和 Harbor，那么在开始离线部署前，你需要先下载这几个工具对应的容器镜像。DevStream 提供了这几个工具对应的镜像列表，并且帮你准备了工具脚本从而更加容易地完成镜像离线工作：

1. [GitLab CE images](../../plugins/gitlab-ce-docker/gitlab-ce-images.txt)
2. [Jenkins images](../../plugins/helm-installer/jenkins/jenkins-images.txt)
3. [Harbor images](../../plugins/helm-installer/harbor/harbor-images.txt)

你可以通过如下命令将镜像列表下载到本地：

```shell
curl -o jenkins-images.txt https://raw.githubusercontent.com/devstream-io/devstream/main/docs/plugins/helm-installer/jenkins/jenkins-images.txt
curl -o harbor-images.txt https://raw.githubusercontent.com/devstream-io/devstream/main/docs/plugins/helm-installer/harbor/harbor-images.txt
curl -o jenkins-images.txt https://raw.githubusercontent.com/devstream-io/devstream/main/docs/plugins/gitlab-ce-docker/gitlab-ce-images.txt
```

可以通过如下命令下载 DevStream 提供的工具脚本，这个脚本可以帮助你快速将这些镜像下载到本地并且上传到私有镜像仓库：

```shell
curl -o image-pull-push.sh https://raw.githubusercontent.com/devstream-io/devstream/main/hack/image-pull-push.sh
chmod +x image-pull-push.sh
```

如果你还没有一个私有镜像仓库，可以参考[这篇文章](../reference/image-registry.zh.md)快速部署一个 Docker Registry。

接下来，你就可以通过下述命令快速完成镜像的下载和上传了：

```shell
# 查看工具脚本的使用方法和注意事项等
./image-pull-push.sh -h # (1)
# 设置镜像仓库地址，按需修改
export IMAGE_REPO_ADDR=registry.devstream.io
# 下载 xxx-images.txt 中所有镜像并保存到本地压缩包中
./image-pull-push.sh -f harbor-images.txt -r ${IMAGE_REPO_ADDR} -s
./image-pull-push.sh -f jenkins-images.txt -r ${IMAGE_REPO_ADDR} -s
./image-pull-push.sh -f gitlab-ce-images.txt -r ${IMAGE_REPO_ADDR} -s
# 从压缩包中 load 镜像并 push 到私有镜像仓库（如果镜像仓库需要登录，则需要先手动执行 docker login）
./image-pull-push.sh -f harbor-images.txt -r ${IMAGE_REPO_ADDR} -l -u
./image-pull-push.sh -f jenkins-images.txt -r ${IMAGE_REPO_ADDR} -l -u
./image-pull-push.sh -f gitlab-ce-images.txt -r ${IMAGE_REPO_ADDR} -l -u
```

1. 强烈建议你先看下本脚本的使用说明和示例

!!! note "注意"

    如果你下载镜像的机器和内部私有镜像仓库之间网络隔离，那么你可以在镜像下载到本地压缩包后，先将该压缩包复制到能够访问镜像仓库的机器上，然后再执行 load 和 push 等操作。

## 3、下载 Helm Chart 包

你可以通过如下命令下载 Harbor 和 Jenkins 的 Helm chart 包：

```shell
helm repo add harbor https://helm.goharbor.io
helm repo update
helm search repo harbor -l
helm pull harbor/harbor --version=1.10.0
```

```shell
helm repo add jenkins https://charts.jenkins.io
helm repo update
helm search repo jenkins -l
helm pull jenkins/jenkins --version=4.2.5
```

执行完上述命令后，你可以在本地看到如下文件：

```shell
$ ls
harbor-1.10.0.tgz jenkins-4.2.5.tgz
```

## 4、准备配置文件

这时候需要联网下载的各种“物料”你就准备好了。接着你可以开始编写 DevStream 的配置文件了：

```yaml title="DevStream Config"
config:
  state:
    backend: local
    options:
      stateFile: devstream.state
vars:
  imageRepo: registry.devstream.io
  gitlabHostname: gitlab.example.com
  jenkinsHostname: jenkins.example.com
  harborHostname: harbor.example.com
  harborURL: http://harbor.example.com
  jenkinsAdminUser: admin
  jenkinsAdminPassword: changeme
  gitlabSSHPort: 30022
  gitlabHttpPort: 30080
  gitlabHttpsPort: 30443

tools:
- name: gitlab-ce-docker
  instanceID: default
  dependsOn: []
  options:
    hostname: [[ gitlabHostname ]]
    gitlabHome: /srv/gitlab
    sshPort: [[ gitlabSSHPort ]]
    httpPort: [[ gitlabHttpPort ]]
    httpsPort: [[ gitlabHttpsPort ]]
    rmDataAfterDelete: false
    imageTag: "rc"
- name: helm-installer
  instanceID: jenkins-001
  dependsOn: []
  options:
    chartPath: "./jenkins-4.2.5.tgz"
    valuesYaml: |
      serviceAccount:
        create: true
        name: jenkins
      controller:
        image: [[ imageRepo ]]/devstreamdev/jenkins
        tag: 2.361.1-jdk11-dtm-0.1
        imagePullPolicy: "IfNotPresent"
        sidecars:
          configAutoReload:
            image: [[ imageRepo ]]/kiwigrid/k8s-sidecar:1.15.0
        adminUser: [[ jenkinsAdminUser ]]
        adminPassword: [[ jenkinsAdminPassword ]]
        ingress:
          enabled: true
          hostName: [[ jenkinsHostname ]]
      enableRawHtmlMarkupFormatter: true
      JCasC:
        defaultConfig: true
- name: helm-installer
  instanceID: harbor-001
  dependsOn: []
  options:
    chartPath: "./harbor-1.10.0.tgz"
    valuesYaml: |
      externalURL: [[ harborURL ]]
      expose:
        type: ingress
        tls:
          enabled: false
        ingress:
          hosts:
            core: [[ harborHostname ]]
      nginx:
        image:
          repository: [[ imageRepo ]]/goharbor/nginx-photon
          tag: v2.5.3
      portal:
        image:
          repository: [[ imageRepo ]]/goharbor/harbor-portal
          tag: v2.5.3
      core:
        image:
          repository: [[ imageRepo ]]/goharbor/harbor-core
          tag: v2.5.3
      jobservice:
        image:
          repository: [[ imageRepo ]]/goharbor/harbor-jobservice
          tag: v2.5.3
      registry:
        registry:
          image:
            repository: [[ imageRepo ]]/goharbor/registry-photon
            tag: v2.5.3
        controller:
          image:
            repository: [[ imageRepo ]]/goharbor/harbor-registryctl
            tag: v2.5.3
      chartmuseum:
          enabled: false
          image:
            repository: [[ imageRepo ]]/goharbor/chartmuseum-photon
            tag: v2.5.3
      trivy:
        enabled: false
        image:
          repository: [[ imageRepo ]]/goharbor/trivy-adapter-photon
          tag: v2.5.3
      notary:
        enabled: false
        server:
          image:
            repository: [[ imageRepo ]]/goharbor/notary-server-photon
            tag: v2.5.3
        signer:
          image:
            repository: [[ imageRepo ]]/goharbor/notary-signer-photon
            tag: v2.5.3
      database:
        internal:
          image:
            repository: [[ imageRepo ]]/goharbor/harbor-db
            tag: v2.5.3
      redis:
        internal:
          image:
            repository: [[ imageRepo ]]/goharbor/redis-photon
            tag: v2.5.3
      exporter:
        image:
          repository: [[ imageRepo ]]/goharbor/harbor-exporter
          tag: v2.5.3
      persistence:
        persistentVolumeClaim:
          registry:
            storageClass: ""
            accessMode: ReadWriteOnce
            size: 5Gi
          jobservice:
            storageClass: ""
            accessMode: ReadWriteOnce
            size: 1Gi
          database:
            storageClass: ""
            accessMode: ReadWriteOnce
            size: 1Gi
          redis:
            storageClass: ""
            accessMode: ReadWriteOnce
            size: 1Gi
```

你可以将这个配置文件保存为 `config.yaml`

## 5、开始部署

现在你可以通过如下命令开始部署 GitLab、Jenkins 和 Harbor 了：

```shell
dtm apply -f config.yaml -y
```

完成部署后，你可以参考[这篇文档](./2-gitlab-jenkins-harbor.zh.md)继续学习如何访问 GitLab、Jenkins 和 Harbor 三个工具。

## 6、环境清理

你可以通过如下命令清理环境：

```shell title="环境清理命令"
dtm delete -f config.yaml -y
```
