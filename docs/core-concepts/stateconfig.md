# State Configuration

In the main config, we can specify which "backend" to use to store DevStream state.

A backend is where to actually store the state file. We support the following types of backend:

- local
- s3

## Local Backend

```yaml
varFile: variables-gitops.yaml

toolFile: tools-gitops.yaml

state:
  backend: local
  options:
    stateFile: devstream.state
```

The `stateFile` under the `options` section is mandatory for local backend.

## S3 Backend

TL;DR: see the config example:

```yaml
varFile: variables-gitops.yaml

toolFile: tools-gitops.yaml

state:
  backend: s3
  options:
    bucket: devstream-remote-state
    region: ap-southeast-1
    key: devstream.state
```

_Note: the `bucket`, `region`, and `key` under the `ptions` section are all mandatory fields for s3 backend._

From release [v0.6.0](https://github.com/devstream-io/devstream/releases/tag/v0.6.0), using AWS S3 to store DevStream state is supported.

More on configuring state [here](./stateconfig.md).

In short, we can use the "backend" keyword to specify where to store the state: either locally or in an S3 bucket. If S3 is used, we need to specify the bucket, region, and the S3 key as well.

### Config File Using S3 Backend

`config.yaml`:

```yaml
varFile: variables-gitops.yaml

toolFile: tools-gitops.yaml

state:
  backend: s3
  options:
    bucket: devstream-test-remote-state
    region: ap-southeast-1
    key: devstream.state
```

`variables-gitops.yaml`:

```yaml
githubUsername: IronCore864
repoName: dtm-test-go
defaultBranch: main

dockerhubUsername: ironcore864

argocdNameSpace: argocd
argocdDeployTimeout: 5m
```

`tools-gitops.yaml`:

```yaml
tools:
- name: github-repo-scaffolding-golang
  instanceID: default
  options:
    owner: [[ githubUsername ]]
    org: ""
    repo: [[ repoName ]]
    branch: [[ defaultBranch ]]
    image_repo: [[ dockerhubUsername ]]/[[ repoName ]]
```

### Using the S3 Backend

Before reading on, now is a good time to check if you have configured your AWS related environment variables correctly or not.

For macOS/Linux users, do:

```shell
export AWS_ACCESS_KEY_ID=ID_HERE
export AWS_SECRET_ACCESS_KEY=SECRET_HERE
export AWS_DEFAULT_REGION=REGION_HERE
```

For more information, see the [official document here](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).

Then let's run `dtm apply`:

```shell
tiexin@mbp ~/work/devstream-io/devstream $ ./dtm apply
2022-05-30 17:07:59 ℹ [INFO]  Apply started.
2022-05-30 17:07:59 ℹ [INFO]  Using dir <.devstream> to store plugins.
2022-05-30 17:07:59 ℹ [INFO]  Using s3 backend. Bucket: devstream-test-remote-state, region: ap-southeast-1, key: devstream.state.
2022-05-30 17:08:00 ℹ [INFO]  Tool (github-repo-scaffolding-golang/default) found in config but doesn't exist in the state, will be created.
Continue? [y/n]
Enter a value (Default is n): y

2022-05-30 17:08:08 ℹ [INFO]  Start executing the plan.
2022-05-30 17:08:08 ℹ [INFO]  Changes count: 1.
2022-05-30 17:08:08 ℹ [INFO]  -------------------- [  Processing progress: 1/1.  ] --------------------
2022-05-30 17:08:08 ℹ [INFO]  Processing: (github-repo-scaffolding-golang/default) -> Create ...
2022-05-30 17:08:12 ℹ [INFO]  The repo dtm-test-go has been created.
2022-05-30 17:08:29 ✔ [SUCCESS]  Tool (github-repo-scaffolding-golang/default) Create done.
2022-05-30 17:08:29 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
2022-05-30 17:08:29 ✔ [SUCCESS]  All plugins applied successfully.
2022-05-30 17:08:29 ✔ [SUCCESS]  Apply finished.
```

As we can see from the output, the S3 backend is used, and it also shows the bucket and key you are using, and in which region this bucket lives.

After `apply`, let's download the state file from S3 and check it out:

```shell
tiexin@mbp ~/work/devstream-io/devstream $ aws s3 cp s3://devstream-test-remote-state/devstream.state .
```

And if we open the downloaded file, we will see something similar to the following:

```yaml
github-repo-scaffolding-golang_default:
  name: github-repo-scaffolding-golang
  instanceid: default
  dependson: []
  options:
    branch: main
    image_repo: ironcore864/dtm-test-go
    org: ""
    owner: IronCore864
    repo: dtm-test-go
  resource:
    org: ""
    outputs:
      org: ""
      owner: IronCore864
      repo: dtm-test-go
      repoURL: https://github.com/IronCore864/dtm-test-go.git
    owner: IronCore864
    repoName: dtm-test-go
```

which is exactly the same as if we were using the local backend to store state.
