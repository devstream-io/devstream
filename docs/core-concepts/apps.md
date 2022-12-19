# Apps

## 1 Concept

An app in DevStream corresponds to a real-world application, and the app represents the whole software development lifecycle of that app, including source code management, code scaffolding, CI/CD (and their pipelines).

Using "App", you can easily create these for an application.

### 1.1 Apps

There are situations where you need to define multiple DevOps tools for an application/microservice. For example, for a web-app typed microservice, you might need the following:

- source code management, code repo scaffolding
- continuous integration (the installation of the DevOps tool, the creation of the CI pipeline)
- continuous deployment (the installation of the DevOps tool, the creation of the CD pipeline)

If you are managing more than one application/microservice (chances are, you will be managing more than one application in the real world), the configuration of DevStream can be quite long, hard to read and hard to manage if you are only using "Tools". 

In order to solve this problem, DevStream provides another concept that is "App". You can easily define all DevOps tools and pipelines for an App with a couple of lines of YAML config, making the config file much easier to read and manage.

In essence, "App" will be converted to "Tool", which you do not have to worry about at all; let DevStream handle that.

## 1.2 pipelineTemplates

pipelineTemplates define CI/CD pipelines, so that they can be referred to and shared by different DevStream Apps, reducing the length of the config file to the next level.

## 2 Config

### 2.1 App

In the config, there is a `apps` section, which is a list, with each element having the following keys:

- name: the name of the app, unique
- spec: application-specific information
- repo: info about the code repository
- repoTemplate: optional, same structure as "repo". If empty, DevStream will create/scaffold a repository from scratch.
- ci: optional, a list of CI pipelines, each element can have the following keys:
    - type: the value can be a `template` or the name of a plugin
    - templateName: optional, if type is `template`, it defines which pipelineTemplate to use
    - vars: optional, variables to be passed to the template. Only works when type is `template`, apparently
    - options: optional
        - if type is the name of a plugin, the options are the options of that plugin
        - if type is `template`, the options here will override the ones in the template. Use full path to override, for example, `options.docker.registry.type`
- cd: like `ci`, but stands for the list of CD pipelines. DevStream will execute CI first before CD

### 2.2 pipelineTemplate

Defined in the `pipelineTemplates` of the config, it's a list, with each element having the following keys:

- name: unique name of the pipelineTemplate, unique
- type: corresponds to a plugin's name
- options: options for that plugin

### 2.3 Local Variables

DevStream has a "var" section in the config, serving as global variables that can be referred to by all Tools and Apps.

Sometimes, however, we'd like to use the same DevOps tool with minor differences. For example, except the name of the project, everything else is different.

In this case, we can define a pipelineTemplate with a local variable, and when referring to it, we can pass different values to it:

```yaml hl_lines="13 15 23 30" title="pipelineTemplate and local variables"
apps:
- name: my-app
  spec:
    language: java
    framework: springboot
    repo: 
      url: https://github.com/testUser/testApp.git
      branch: main
    ci:
    - type: github-actions # use a plugin directly without defining pipelineTemplates
    cd:
    - type: template # use a pipelineTemplate
      templateName: my-cd-template # corresponds to the name of the pipelineTemplate
      vars:
        appName: my-app # a local variable passed to the pipelineTemplate

pipelineTemplates:
cd:
- name: my-cd-template
  type: argocdapp
  options:
    app:
      name: [[ appName ]] # a local variable, passed to when referring to the template
      namespace: argocd
    destination:
      server: https://kubernetes.default.svc
      namespace: default
    source:
      valuefile: values.yaml
      path: charts/[[ appName ]]
```

## 3 A Demo of the Whole Config

A whole config for an App:

```yaml
apps:
- name: testApp # name of the app
  spec: # app-specific info
    language: java # programming language of the app
    framework: springboot # framework of the app
  repo: # repository-related info for the app
    url: https://github.com/testUser/testApp.git
    branch: main
  repoTemplate: # optional, used for repository bootstrapping/scaffolding. If not empty, a repo will be created with scaffolding code
    url: https://github.com/devstream-io/dtm-repo-scaffolding-java-springboot.git
    vars:
      imageRepoOwner: repoOwner # variables used for repoTemplate
  ci: # CI pipelines, here we use github-actions
  - type: github-actions
- name: testApp2
  spec:
    language: go
    framework: gin
  repo: # repository-related info for the app
    owner: test_user
    type: github
    branch: main
  repoTemplate: # optional, used for repository bootstrapping/scaffolding. If not empty, a repo will be created with scaffolding code
    org: devstream-io
    name: dtm-repo-scaffolding-java-springboot
    type: github
  ci: # CI pipelines, here we use github-actions
  - type: github-actions
    options:
      imageRepo:
        owner: repoOwner # override the plugin's options. Must use full YAML path.
  cd: # CD pipelines, here we use argocd
  - type: argocdapp
```

If we apply this config, DevStream will create two repositories in GitHub, with scaffolding code provided by DevStream [SpringBoot](https://github.com/devstream-io/dtm-repo-scaffolding-java-springboot.git). App `testApp` will trigger CI in GitHub Actions upon each commit, and App `testApp2` will trigger build/push in GitHub Actions upon commit, and deploy using Argo CD.

### repo/repoTemplate Config

The `repo` and `repoTemplate` in the Config represent a code repository. You can define it with a single URL or a few key/values:

!!! note "two ways to configure code repo"

    === "using a single URL"
    
        ```yaml title=""
          repo:
            url: git@gitlab.example.com:root/myapps.git # repo URL, supports both git and https
            apiURL: https://gitlab.example.com # not mandatory, if using gitlab and the URL protocol is git, here can be the GitLab API URL
            branch: "" # not mandatory, defaults to main for GitHub and master for GitLab
        ```
    
        This example shows that we use GitLab `git@gitlab.example.com:root/myapps.git` for code clone, and DevStream uses `https://gitlab.example.com` to access GitLab API. Default branch is master.
    
    === "using detailed key/value config for the repo"
    
        ```yaml title=""
          repo:
            org: "" # only mandatory for GitHub organization
            owner："test_user" # if the repo belongs to a person. If the repo belongs to an org, use the org above.
            name: "" # optional, defaults to the name of the app
            baseURL:  https://gitlab.example.com # optional. If GitLab, here we can put the GitLab domain.
            branch: master # not mandatory, defaults to main for GitHub and master for GitLab
            type:  gitlab # mandatory, either gitlab or github
        ```
    
        This example shows that we use GitLab `https://gitlab.example.com`, repo name is the app name, belongs to owner `test_user`, with the default branch being "master".

### CI Config

The `CI` section in the config supports 4 types at the moment: `github-actions`/`gitlab-ci`/`jenkins-pipeline`/`template`.

`template` means to use a pipelineTemplate; and the other three types correspond to GitHub Actions, GitLab CI, and Jenkins, respectively.

Detailed config:

```yaml
  ci:
  - type: jenkins-pipieline # type of the CI
    options: # options for CI. If empty, CI will only run unit test.
      jenkins: # config for jenkins
        url: jenkins.exmaple.com # jenkins URL
        user: admin # jenkins user
      imageRepo: # docker image repo to be pushed to. If set, Ci will push the image after build.
        url: http://harbor.example.com # image repo URL. Defaults to dockerhub.
        owner: admin # image repo owner
      dingTalk: # dingtalk notification settings. If set, CI result will be pushed to dingtalk.
        name: dingTalk
        webhook: https://oapi.dingtalk.com/robot/send?access_token=changemeByConfig # callback URL for dingtalk.
        securityType: SECRET # use secret to encrypt dingtalk message
        securityValue: SECRETDATA # dingtalk secret encryption data
      sonarqube: # sonarqube config. If set, CI will test and execute sonarqube scan.
        url: http://sonar.example.com # sonarqube URL
        token: YOUR_SONAR_TOKEN # soanrqube token
        name: sonar_test
```

The config above will trigger unit test and sonarqube code scan upon commit, then a Docker image will be built and pushed to dockerhub, and the result of the CI will be pushed to dingtalk.

If the same pipeline is required for multiple apps, the config can be long and redundant. So, DevStream provides the `template` type to share similar settings for diffrent Apps. Detailed example:

```yaml
apps:
- name: javaProject1
  spec:
    language: java
    framework: springboot
  repo:
    owner: testUser
    type: github
  repoTemplate:
    url: https://github.com/devstream-io/dtm-repo-scaffolding-java-springboot.git
  ci:
  - type: template # use a pipelineTemplate
    templateName: ci-pipeline # name of the pipelineTemplate
    vars:
      dingdingAccessToken: tokenForProject1 # variables for the pipelineTemplate
      dingdingSecretValue: secretValProject1
- name: javaProject2
  spec:
    language: java
    framework: springboot
  repo:
    owner: testUser
    type: github
  repoTemplate:
    url: https://github.com/devstream-io/dtm-repo-scaffolding-java-springboot.git
  ci:
  - type: template # use a pipelineTemplate
    templateName: ci-pipeline # name of the pipelineTemplate
    vars:
      dingdingAccessToken: tokenForProject2 # variables for the pipelineTemplate
      dingdingSecretValue: secretValProject2

pipelineTemplates: # CI/CD pipeline templates
- name: ci-pipeline # name of the pipelineTemplate
  type: jenkins-pipeline # type, supports jenkins-pipeline，github-actions and gitlab-ci at the moment
  options: # options, same as CI options
    jenkins:
      url: jenkins.exmaple.com
      user: admin
    imageRepo:
      url: http://harbor.example.com
      owner: admin
    dingTalk:
      name: dingTalk
      webhook: https://oapi.dingtalk.com/robot/send?access_token=[[ dingdingAccessToken ]] # local variable, passed to when referring to this template
      securityType: SECRET
      securityValue: [[ dingdingSecretValue ]] # local variable, passed to when referring to this template
    sonarqube:
        url: http://sonar.example.com
        token: sonar_token
        name: sonar_test

```

If we apply the above config, we will create two Jenkins pipelines for two apps, with the only difference being that the dingtalk notification will be sent to different groups.

### CD Config

At the moment, CD only supports `argocdapp`. Argo CD itself can be deployed with a Tool, and `argocdapp` is responsible for deploying the app in a Kubernetes cluster.

Detailed config example:

```yaml
cd:
- type: argocdapp
  options:
    app:
      name: hello # argocd app name
      namespace: argocd # argocd namespace
    destination:
      server: https://kubernetes.default.svc # Kubernetes cluster
      namespace: default # which namespace to deploy the app
    source:
      valuefile: values.yaml # helm values file
      path: charts/go-hello-http # helm chart path
```
