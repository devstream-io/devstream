# 用 DTM Apps 实现基于 GitHub，Argo CD 和 GitHub Actions 的 CICD 流程

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
- name: helm-installer
  instanceID: argocd

apps:
- name: myapp1
  spec:
    language: python
    framework: django
  repo:
    url: github.com/[[ GITHUB_USER ]]/myapp1
    token: [[ env GITHUB_TOKEN ]]
  repoTemplate:
    url: github.com/devstream-io/dtm-repo-scaffolding-python-flask
  ci:
  - type: github-actions
    options:
      imageRepo:
        user: [[ DOCKERHUB_USER ]]
        password: [[ env IMAGE_REPO_PASSWORD ]]
  cd:
  - type: argocdapp
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

<script id="asciicast-z1XlVm9kGjzArV9aNERD7acfH" src="https://asciinema.org/a/z1XlVm9kGjzArV9aNERD7acfH.js" async></script>
