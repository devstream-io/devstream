# jenkins-pipeline 插件离线配置

## 功能

目前的 `jenkins-pipeline` 插件配置过程中会连接外部网络来安装插件和配置 `pipeline` 的共享库，为了支持在离线环境中使用 `jenkins-pipeline` 来创建和管理 `jenkins` 流水线，我们在 `jenkins-pipeline` 插件中提供了一个选项 `offline` 用于表示需要在离线环境下配置 `jenkins-pipeline`。

## 具体实现

### 插件

流水线正常运行需要依赖 `jenkins` 的几个插件。在离线环境中因为无法下载外部环境的插件，需要在 `jenkins` 预先安装好插件或者使用已经有插件的 `jenkins` 镜像，devstream 官方也提供了已预先安装好所有依赖的镜像 `devstreamdev/jenkins:2.361.1-jdk11-dtm-0.1`。

依赖的插件如下：
| 插件名       | 作用           | 备注                       |
|--------------|----------------|----------------------------|
| dingding-notifications    | 用于发送钉钉通知          | 如果需要使用钉钉通知，则需要安装此插件                 |
| github-branch-source      | jenkins github 插件       | 如果代码仓库使用的是 github，则必须安装此插件          |
| gitlab-plugin             | jenkins gitlab 插件       | 如果代码仓库使用的是 gitlab，则必须安装此插件          |
| sonar                     | sonar sanner 代码扫描插件 | 如果需要使用 sonarqube，则需要安装此插件               |
| kubernetes                | jenkins kubernetes 插件   | 用于 jenkins runner， 必须安装此插件                   |
| git                       | jenkins git 插件          | 用于代码的克隆和权限认证， 必须安装此插件              |
| configuration-as-code     | jenkins 配置即代码插件    | 用于 devstream 配置 jenkins 的全局配置， 必须安装此插件|

### 共享库

在 `jenkins` 中 devstream 默认使用共享库来管理 jenkins 流水线的共用代码，在离线环境中因为无法连接共享库，所以 devstream 提供了单独的 `Jenkinfile` 配置。


## app 示例配置

```yaml
apps:
- name: test
  spec:
    language: java
    framework: springboot
  repo:
    url: gitlab.com/root/test.git
    branch: main
  repoTemplate:
    url: https://github.com/devstream-io/dtm-repo-scaffolding-java-springboot.git
  ci:
  - type: template
    templateName: ci-pipeline

pipelineTemplates:
- name: ci-pipeline
  type: jenkins-pipeline
  options:
    branch: main
    jenkins:
      url: jenkins.com
      user: admin
      offline: true # 在此处设置 offline 为 true， 即开启该 jenkins-pipeline 的离线模式
    imageRepo:
      user: repoUser
```

使用该配置可得到以下输出：

```text
2022-12-02 19:51:52 ℹ [INFO]  Apply started.
2022-12-02 19:51:52 ℹ [INFO]  Using local backend. State file: devstream-app.state.
2022-12-02 19:51:52 ℹ [INFO]  Tool (repo-scaffolding/test) found in config but doesn't exist in the state, will be created.
2022-12-02 19:51:52 ℹ [INFO]  Tool (jenkins-pipeline/test) found in config but doesn't exist in the state, will be created.
2022-12-02 19:51:52 ℹ [INFO]  Start executing the plan.
2022-12-02 19:51:52 ℹ [INFO]  Changes count: 2.
2022-12-02 19:51:52 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
2022-12-02 19:51:52 ℹ [INFO]  Processing: (repo-scaffolding/test) -> Create ...
2022-12-02 19:51:52 ℹ [INFO]  github start to download repoTemplate...
2022-12-02 19:51:57 ✔ [SUCCESS]  Tool (repo-scaffolding/test) Create done.
2022-12-02 19:51:57 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
2022-12-02 19:51:57 ℹ [INFO]  Processing: (jenkins-pipeline/test) -> Create ...
2022-12-02 19:51:58 ℹ [INFO]  jenkins plugin imageRepo start config...
2022-12-02 19:51:58 ⚠ [WARN]  jenkins gitlab ssh key not config, private repo can't be clone
2022-12-02 19:52:00 ℹ [INFO]  jenkins start config casc...
2022-12-02 19:52:07 ✔ [SUCCESS]  Tool (jenkins-pipeline/test) Create done.
2022-12-02 19:52:07 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
2022-12-02 19:52:07 ✔ [SUCCESS]  All plugins applied successfully.
2022-12-02 19:52:07 ✔ [SUCCESS]  Apply finished.
```