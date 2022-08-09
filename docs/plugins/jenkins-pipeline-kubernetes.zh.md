# jenkins-pipeline-kubernetes 插件

这个插件在已有的 Jenkins 上建立 Jenkins job, 将 GitHub 作为 SCM。

步骤：

1. 按需修改配置项，其中 `githubRepoUrl` 为 GitHub 仓库地址，应预先建立一个 GitHub 仓库，并创建一个名为 "Jenkinsfile" 的文件放至仓库根目录。
2. 设置环境变量
    - `GITHUB_TOKEN`
    - `JENKINS_PASSWORD`

## 用例

```yaml
--8<-- "jenkins-pipeline-kubernetes.yaml"
```

## 和 `jenkins` 插件一起使用

这个插件可以和 `jenkins` 插件一起使用，[`jenkins` 插件文档](./jenkins.zh.md)。

即在安装完 `Jenkins` 后，再建立 `Jenkins` job。

首先根据 `dependsOn` 设定插件依赖，再根据 `${{jenkins.default.outputs.jenkinsURL}}` 和 `${{jenkins.default.outputs.jenkinsPasswordOfAdmin}}` 设置 Jenkins 的 URL 和 admin 密码。

注意：如果你的 Kubernetes 集群是 K8s in docker 模式，请自行设置网络，确保 `jenkins` 插件中 `NodePort` 设置的端口在宿主机内能访问。

```yaml
---
tools:
  # name of the tool
  - name: jenkins
    # id of the tool instance
    instanceID: default
    # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
    dependsOn: [ ]
    # options for the plugin
    options:
      # if true, the plugin will use hostpath to create a pv named `jenkins-pv`
      # and you should create the volumes directory manually, see plugin doc for details.
      test_env: false
      # need to create the namespace or not, default: false
      create_namespace: false
      # Helm repo information
      repo:
        # name of the Helm repo
        name: jenkins
        # url of the Helm repo
        url: https://charts.jenkins.io
      # Helm chart information
      chart:
        # name of the chart
        chart_name: jenkins/jenkins
        # release name of the chart
        release_name: dev
        # k8s namespace where jenkins will be installed
        namespace: jenkins
        # whether to wait for the release to be deployed or not
        wait: true
        # the time to wait for any individual Kubernetes operation (like Jobs for hooks). This defaults to 5m0s
        timeout: 5m
        # whether to perform a CRD upgrade during installation
        upgradeCRDs: true
        # custom configuration. You can refer to [Jenkins values.yaml](https://github.com/jenkinsci/helm-charts/blob/main/charts/jenkins/values.yaml)
        values_yaml: |
          persistence:
            # for prod env: the existent storageClass, please change it
            # for test env: just ignore it, but don't remove it
            storageClass: jenkins-pv
          serviceAccount:
            create: false
            name: jenkins
          controller:
            serviceType: NodePort
            nodePort: 32000
            additionalPlugins:
              # install "GitHub Pull Request Builder" plugin, see https://plugins.jenkins.io/ghprb/ for more details
              - ghprb
              # install "OWASP Markup Formatter" plugin, see https://plugins.jenkins.io/antisamy-markup-formatter/ for more details
              - antisamy-markup-formatter
            # Enable HTML parsing using OWASP Markup Formatter Plugin (antisamy-markup-formatter), useful with ghprb plugin.
            enableRawHtmlMarkupFormatter: true
            # Jenkins Configuraction as Code, refer to https://plugins.jenkins.io/configuration-as-code/ for more details
            # notice: All configuration files that are discovered MUST be supplementary. They cannot overwrite each other's configuration values. This creates a conflict and raises a ConfiguratorException.
            JCasC:
              defaultConfig: true
              # each key-value in configScripts will be added to the ${JENKINS_HOME}/casc_configs/ directory as a file.
              configScripts:
                # this will create a file named "safe_html.yaml" in the ${JENKINS_HOME}/casc_configs/ directory.
                # it is used to configure the "Safe HTML" plugin.
                # filename must meet RFC 1123, see https://tools.ietf.org/html/rfc1123 for more details
  - name: jenkins-pipeline-kubernetes
    # id of the tool instance
    instanceID: default
    # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
    dependsOn: [ "jenkins.default" ]
    # options for the plugin
    options:
      jenkins:
        # jenkinsUrl, format: hostname:port
        url: ${{jenkins.default.outputs.jenkinsURL}}
        # jenkins user name, default: admin
        user: admin
        # jenkins password, you have 3 options to set the password:
        # 1. use outputs of the `jenkins` plugin, see docs for more details
        # 2. set the `JENKINS_PASSWORD` environment variable
        # 3. fill in the password in this field(not recommended)
        # if all set, devstream will read the password from the config file or outputs from jenkins plugin first, then env var.
        password: ${{jenkins.default.outputs.jenkinsPasswordOfAdmin}}
        # jenkins job name, mandatory
        jobName:
        # path to the pipeline file, relative to the git repo root directory. default: Jenkinsfile
        pipelineScriptPath: Jenkinsfile
      # github repo url where the pipeline script is located. mandatory
      githubRepoUrl: https://github.com/YOUR_GITHUB_ACCOUNT/YOUR_TEST_PROJECT_NAME.git
```

