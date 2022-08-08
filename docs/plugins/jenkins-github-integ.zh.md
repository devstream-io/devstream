# jenkins-github-integ 插件

本插件可以在基于 helm 安装的 Jenkins 上集成 GitHub。主要工作如下：

- 安装 [GitHub Pull Request Builder](https://plugins.jenkins.io/ghprb/) 插件、 [OWASP Markup Formatter](https://plugins.jenkins.io/antisamy-markup-formatter/) 插件。
- 配置 Jenkins 的 Pull Request Builder 插件，自动为 GitHub repo 创建 WebHook。
- 创建一个基于 GitHub pull request 触发的 Jenkins job。

建议与 [`jenkins` 插件](./jenkins.zh.md)配置使用。

## ngrok 使用
如果处在测试环境，无公网 IP，为了让 GitHub 的 WebHook 能访问到 Jenkins，需要将本机的 Jenkins 暴露到公网。可以选取 [ngrok](https://ngrok.com/) 作为内网穿透工具。

假设你的 Jenkins 的端口为 8080，安装 ngrok, 执行 `ngrok http 8080` ，复制 "Forwarding" 行的公网访问地址，粘贴到 dtm 配置文件中的 `options.jenkins.urlOverride` 字段即可。

例：

```yaml
tools:
  - name: jenkins-github-integ
    instanceID: default
    dependsOn: [ "jenkins.default" ]
    options:
      jenkins:
        url: http://localhost:8080
        urlOverride: https://891e-125-111-206-162xxx.ap.ngrok.io
    ...
    ...
```


## 用例

```yaml

--8<-- "jenkins-github-integ.yaml"

```

## 和 `jenkins` 插件一起使用

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
  - name: jenkins-github-integ
    # id of the tool instance
    instanceID: default
    # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
    dependsOn: [ "jenkins.default" ]
    # options for the plugin
    options:
      jenkins:
        # jenkinsUrl, format: hostname:port
        url: ${{jenkins.default.outputs.jenkinsURL}}
        # override the jenkins url to expose the jenkins to the GitHub webhook, can be empty.
        urlOverride:
        # jenkins user name, default: admin
        user: admin
        # jenkins password, you have 3 options to set the password:
        # 1. use outputs of the `jenkins` plugin, see docs for more details
        # 2. set the `JENKINS_PASSWORD` environment variable
        # 3. fill in the password in this field(not recommended)
        # if all set, devstream will read the password from the config file first.
        password: ${{jenkins.default.outputs.jenkinsPasswordOfAdmin}}
        # jenkins job name, mandatory
        jobName: 
        # path to the pipeline file, relative to the git repo root directory. default: Jenkinsfile-pr
        pipelineScriptPath: Jenkinsfile-pr
      helm:
        # namespace of the jenkins, default: jenkins
        namespace: jenkins
        # release name of the jenkins helm chart, mandatory
        releaseName: 
      # GitHub repo where to put the pipeline script and project. mandatory
      githubRepoUrl: https://github.com/YOUR_GITHUB_ACCOUNT/YOUR_TEST_PROJECT_NAME
      adminList:
        - YOUR_GITHUB_USERNAME
```

