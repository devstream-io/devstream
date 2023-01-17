# GitHub + GitHub Actions with DTM Apps

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

Update the "YOUR_GITHUB_USER" and "YOUR_DOCKERHUB_USER" in the above file accordingly.

---

## Run

First, initialize:

```bash
# this downloads the required plugins, according to the config file, automatically.
dtm init -f config.yaml
```

<script id="asciicast-EMzx8lzZq5AJpAMoJY23fh8A3" src="https://asciinema.org/a/EMzx8lzZq5AJpAMoJY23fh8A3.js" async></script>

Then we apply it by running:

```bash
dtm apply -f config.yaml -y
```

<script id="asciicast-z1XlVm9kGjzArV9aNERD7acfH" src="https://asciinema.org/a/z1XlVm9kGjzArV9aNERD7acfH.js" async></script>
