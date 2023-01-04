# 用 DTM Tools 实现基于 GitHub，Argo CD 和 GitHub Actions 的 CICD 流程

DevStream 抽象了2个概念：[Tools](../../core-concepts/tools.md) 和 [Apps](../../core-concepts/apps.md)。

在[前一个用户场景](./2-github-dtm-apps.md) 里介绍了 "Apps"，你可以用 "Tools" 实现一样的效果。具体方法如下：

## 环境变量

你需要先设置如下环境变量：

```bash
export GITHUB_TOKEN="YOUR_GITHUB_TOKEN_HERE"
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
  GITHUB_USER: YOUR_GITHUB_USER
  DOCKERHUB_USER: YOUR_DOCKERHUB_USER

tools:
- name: repo-scaffolding
  instanceID: myapp1
  options:
    destinationRepo:
      owner: [[ GITHUB_USER ]]
      name: myapp1
      branch: main
      scmType: github
      token: [[ env GITHUB_TOKEN ]]
    sourceRepo:
      org: devstream-io
      name: dtm-repo-scaffolding-python-flask
      scmType: github
- name: github-actions
  instanceID: flask
  dependsOn: [ repo-scaffolding.myapp1 ]
  options:
    scm:
      owner: [[ GITHUB_USER ]]
      name:  myapp1
      scmType: github
      token: [[ env GITHUB_TOKEN ]]
    pipeline:
      configLocation: https://raw.githubusercontent.com/devstream-io/dtm-pipeline-templates/main/github-actions/workflows/main.yml
      language:
        name: python
        framework: flask
      imageRepo:
        user: [[ DOCKERHUB_USER ]]
        password: [[ env IMAGE_REPO_PASSWORD ]]
- name: helm-installer
  instanceID: argocd
- name: argocdapp
  instanceID: default
  dependsOn: [ "helm-installer.argocd", "github-actions.flask" ]
  options:
    app:
      name: myapp1
      namespace: argocd
    destination:
      server: https://kubernetes.default.svc
      namespace: default
    source:
      valuefile: values.yaml
      path: helm/myapp1
      repoURL: ${{repo-scaffolding.myapp1.outputs.repoURL}}
      token: [[ env GITHUB_TOKEN ]]
    imageRepo:
      user: [[ DOCKERHUB_USER ]]
```

你需要相应更新上述配置文件里的 "YOUR_GITHUB_USER" 和 "YOUR_DOCKERHUB_USER"。

---

## 运行

首先需要初始化：

```bash
# this downloads the required plugins, according to the config file, automatically.
dtm init -f config.yaml
```

<script id="asciicast-EMzx8lzZq5AJpAMoJY23fh8A3" src="https://asciinema.org/a/EMzx8lzZq5AJpAMoJY23fh8A3.js" async></script>

然后运行如下命令让配置生效：

```bash
dtm apply -f config.yaml -y
```

(省略了动图和视频等)
