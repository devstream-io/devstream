# dtm-repo-scaffolding-python-flask

This repo contains templates used by DevStream plugin "repo-scaffolding" (thereafter: the plugin).

This repo isn't intended to be used directly without DevStream. It should only be consumed by the plugin automatically.

The plugin (together with this repo of templates) can create a repo in GitHub and set up the project layout and initialize the reop with necessary files that are typical for a Go web app. The followings can be created automatically:

- a Python web app example using the Flask framework
- directory structure, following Flask/Python best practice
- `.gitignore`, suggested by Flask
- Dockerfile with Python alpline
- a simplified Helm chart with Deployment and Service

## Usage

- Render all files using go template whose name end with `.tpl` suffix.
- Files whose name don't end with `.tpl` extension don't need to be rendered.
- subdirectory "helm/**_app_name_**" (the **_app_name_** part) should be rendered with `AppName`
- subdicrectory "cmd/**_app_name_**" (the **_app_name_** part) should be rendered with `AppName`

Example of required parameters to render these templates:

```yaml
AppName: hello
Repo:
  Owner: ironcore864
  Name: hello
imageRepo: ironcore864/hello # dockerhub
```

## Where does this repo come from?

`dtm-repo-scaffolding-python-flask` is synced from https://github.com/devstream-io/devstream/blob/main/staging/dtm-repo-scaffolding-python-flask. 
Code changes are made in that location, merged into `devstream-io/devstream` and later synced here.
