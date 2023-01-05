# 用 DTM Tools 实现基于 Gitlab，Argo CD 和 Gitlab CI 的 CICD 流程

DevStream 抽象了2个概念：[Tools](../../core-concepts/tools.md) 和 [Apps](../../core-concepts/apps.md)。

在[前一个用户场景](./4-gitlab-dtm-apps.zh.md) 里介绍了 "Apps"，你可以用 "Tools" 实现一样的效果。具体方法如下：

## 配置环境变量

你需要配置以下两个环境变量：

```bash
export GITLAB_TOKEN="YOUR_GITLAB_TOKEN"
export IMAGE_REPO_PASSWORD="YOUR_DOCKERHUB_TOKEN_HERE"
```

---

## 配置文件

```yaml
config:
  state:
    backend: local
    options:
      stateFile: devstream.state

vars:
  gitlabUser: YOUR_GITLAB_USERNAME
  dockerUser: YOUR_DOCKERHUB_USERNAME
  app: testapp

tools:
- name: helm-installer
  instanceID: argocd
- name: repo-scaffolding
  instanceID: myapp
  options:
    destinationRepo:
      owner: [[ gitlabUser ]]
      name: [[ app ]]
      branch: master
      scmType: gitlab
      # set env GITLAB_TOKEN
      token: [[ env GITLAB_TOKEN ]]
    sourceRepo:
      org: devstream-io
      name: dtm-repo-scaffolding-python-flask
      scmType: github
- name: gitlab-ci
  instanceID: flask
  dependsOn: [ "repo-scaffolding.myapp" ]
  options:
    scm:
      owner: [[ gitlabUser ]]
      name: [[ app ]]
      branch: master
      scmType: gitlab
      token: [[ env GITLAB_TOKEN ]]
    pipeline:
      language:
        framework: flask
        name: python
      imageRepo:
        user: [[ dockerUser ]]
        # set env IMAGE_REPO_PASSWORD
        password: [[ env IMAGE_REPO_PASSWORD ]]
- name: argocdapp
  instanceID: default
  dependsOn: [ "helm-installer.argocd", "gitlab-ci.flask" ]
  options:
    app:
      name: [[ app ]]
      namespace: argocd
    destination:
      server: https://kubernetes.default.svc
      namespace: default
    source:
      valuefile: values.yaml
      path: helm/[[ app ]]
      repoURL: ${{repo-scaffolding.myapp.outputs.repoURL}}
      token: [[ env GITLAB_TOKEN ]]
    imageRepo:
      user: [[ dockerUser ]]
      password: [[ env IMAGE_REPO_PASSWORD ]]
```

需要相应更新上述配置文件里的 "YOUR_GITLAB_USERNAME" 和 "YOUR_DOCKERHUB_USERNAME"。

**注意：**

上述配置运行的前提是你的 `GitLab` 已经配置好了共享 runner 来执行 `ci` 操作，如果你想要自动创建一个 runner， 你可以参考 [gitlab插件文档](../../plugins/github-actions.zh.md) 来配置由 `DevStream` 来创建一个项目 runner。

---

## 执行

首先需要基于配置文件来初始化插件：

```bash
# this downloads the required plugins, according to the config file, automatically.
dtm init -f config.yaml
```

然后运行如下命令让配置生效：

```bash
dtm apply -f config.yaml -y
```

现在我们就可以看到 `GitLab` 中已经创建了对应的仓库并且 `CI` 把生成的项目镜像推送到了 `Dockerhub` 仓库。

<figure markdown>
  ![GitLab CI ](./gitlab-tools/gitlab-tools-ci-page.png){ width="1000" }
  <figcaption>GitLab CI</figcaption>
</figure>

在你的 `Kubernetes` 集群中已经创建了对应的应用。

```bash
$ kubectl get application -n argocd
NAME      SYNC STATUS   HEALTH STATUS
testapp   Synced        Healthy
$ kubectl get deploy -n default
NAME      READY   UP-TO-DATE   AVAILABLE   AGE
testapp   1/1     1            1           4m27s
$ kubectl get pods -n default
NAME                       READY   STATUS    RESTARTS   AGE
testapp-5f9c75b4f6-57d9p   1/1     Running   0          3m48s
```