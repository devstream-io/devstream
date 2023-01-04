# GitHub + GitHub Actions with DTM Tools

DevStream has two abstractions: [Tools](../../core-concepts/tools.md) and [Apps](../../core-concepts/apps.md).

[The previous use case](./2-github-dtm-apps.md) uses "Apps". We can also achieve the same result with "Tools", and here's how:

## ENV Vars

The following environment variables are required for this to work:

```bash
export GITHUB_TOKEN="YOUR_GITHUB_TOKEN_HERE"
export IMAGE_REPO_PASSWORD="YOUR_DOCKERHUB_TOKEN_HERE"
```

---

## Config File

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

Update the "YOUR_GITHUB_USER" and "YOUR_DOCKERHUB_USER" in the above file accordingly.

---

## Run

First, initialize:

```bash
# this downloads the required plugins, according to the config file, automatically.
dtm init -f config.yaml
```

Then we apply it by running:

```bash
dtm apply -f config.yaml -y
```

(Screenshot/video omitted.)
