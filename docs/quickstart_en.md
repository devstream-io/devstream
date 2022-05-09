# Quick Start

> English | [中文](./quickstart_zh.md)

## 1 Download DevStream (`dtm`)

Download the appropriate `dtm` version for your platform from [DevStream Releases](https://github.com/devstream-io/devstream/releases).

> Remember to rename the binary file to `dtm` so that it's easier to use. For example: `mv dtm-darwin-arm64 dtm`.

> Once downloaded, you can run the binary from anywhere. Ideally, you want to put it in a place that is in your PATH (e.g., `/usr/local/bin`).

## 2 Prepare a Config File

Copy the [examples/quickstart.yaml](../examples/quickstart.yaml) and [examples/tools-quickstart.yaml](../examples/tools-quickstart.yaml) to your working directory and rename `quickstart.yaml` to `config.yaml`:

```shell
cp examples/quickstart.yaml config.yaml
cp examples/tools-quickstart.yaml tools-quickstart.yaml
```

Then modify the file accordingly.

For example, my GitHub username is "IronCore864", and my Dockerhub username is "ironcore864", then I can run:

```shell
sed -i.bak "s/YOUR_GITHUB_USERNAME_CASE_SENSITIVE/IronCore864/g" tools-quickstart.yaml

sed -i.bak "s/YOUR_DOCKER_USERNAME/ironcore864/g" tools-quickstart.yaml
```

> This config file uses two plugins, one will create a GitHub repository and bootstrap it into a Golang web app, and the other will create GitHub Actions workflow for it.

The two plugins [require an environment variable](https://www.devstream.io/docs/plugins/github-repo-scaffolding-golang) to work, so let's set it:

```shell
export GITHUB_TOKEN="YOUR_GITHUB_TOKEN_HERE"
```

If you don't know how to create a GitHub token, check out [the official document here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token).

## 3 Initialize

Run:

```shell
dtm init -f config.yaml
```

and you should see similar output to:

```
2022-03-04 12:08:06 ℹ [INFO]  Initialize started.
2022-03-04 12:08:06 ℹ [INFO]  Using dir <.devstream> to store plugins.
2022-03-04 12:08:11 ℹ [INFO]  Downloading: [github-repo-scaffolding-golang-darwin-arm64_0.2.0.so] ...
 13.52 MiB / 13.52 MiB [=================================] 100.00% 3.14 MiB/s 4s
2022-03-04 12:08:15 ✔ [SUCCESS]  [github-repo-scaffolding-golang-darwin-arm64_0.2.0.so] download succeeded.
2022-03-04 12:08:17 ℹ [INFO]  Downloading: [githubactions-golang-darwin-arm64_0.2.0.so] ...
 16.05 MiB / 16.05 MiB [=================================] 100.00% 5.41 MiB/s 2s
2022-03-04 12:08:20 ✔ [SUCCESS]  [githubactions-golang-darwin-arm64_0.2.0.so] download succeeded.
2022-03-04 12:08:20 ✔ [SUCCESS]  Initialize finished.
```

This step verifies the MD5 sum of your dtm binary, downloads the required plugins according to the config file, and verifies the plugins' MD5 sums as well.

Note: if your dtm binary's MD5 sum doesn't match the MD5 sum from our release page, dtm init will stop. If your local dtm MD5 differs, it indicates that you built the binary yourself (for developing purposes, for example). Due to the nature of the Go plugin, dtm must be built together with the corresponding plugins. So, if you are building dtm, you should also build the plugins as well, in which case, you do not need to run dtm init to download the plugins.

## 4 Apply

Run:

```shell
dtm apply -f config.yaml
```

and confirm to continue, then you should see similar output to:

```
2022-03-04 12:08:54 ℹ [INFO]  Apply started.
2022-03-04 12:08:54 ℹ [INFO]  Using dir <.devstream> to store plugins.
2022-03-04 12:08:54 ℹ [INFO]  Tool < go-webapp-repo > found in config but doesn't exist in the state, will be created.
2022-03-04 12:08:54 ℹ [INFO]  Tool < golang-demo-actions > found in config but doesn't exist in the state, will be created.
Continue? [y/n]
Enter a value (Default is n): y

2022-03-04 12:08:57 ℹ [INFO]  Start executing the plan.
2022-03-04 12:08:57 ℹ [INFO]  Changes count: 2.
2022-03-04 12:08:57 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
2022-03-04 12:08:57 ℹ [INFO]  Processing: go-webapp-repo -> Create ...
2022-03-04 12:09:04 ℹ [INFO]  Repo created.
2022-03-04 12:09:22 ✔ [SUCCESS]  Plugin go-webapp-repo Create done.
2022-03-04 12:09:22 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
2022-03-04 12:09:22 ℹ [INFO]  Processing: golang-demo-actions -> Create ...
2022-03-04 12:09:23 ℹ [INFO]  Language is: go-1.17.
2022-03-04 12:09:23 ℹ [INFO]  Creating GitHub Actions workflow pr-builder.yml ...
2022-03-04 12:09:24 ✔ [SUCCESS]  Github Actions workflow pr-builder.yml created.
2022-03-04 12:09:25 ℹ [INFO]  Creating GitHub Actions workflow main-builder.yml ...
2022-03-04 12:09:26 ✔ [SUCCESS]  Github Actions workflow main-builder.yml created.
2022-03-04 12:09:26 ✔ [SUCCESS]  Plugin golang-demo-actions Create done.
2022-03-04 12:09:26 ✔ [SUCCESS]  All plugins applied successfully.
2022-03-04 12:09:26 ✔ [SUCCESS]  Apply finished.
```
## 5 Check the Results

Go to your GitHub account, and we can see a new repo named "go-webapp-devstream-demo" has been created; there are some Golang web app scaffolding lying around already, and the GitHub Actions for building the app is also ready. Hooray!

## 6 Clean Up

Run:

```shell
dtm destroy
```

and you should see similar output:

```
2022-03-04 12:10:36 ℹ [INFO]  Destroy started.
2022-03-04 12:10:36 ℹ [INFO]  Change added: go-webapp-repo_github-repo-scaffolding-golang -> Delete
2022-03-04 12:10:36 ℹ [INFO]  Change added: golang-demo-actions_githubactions-golang -> Delete
Continue? [y/n]
Enter a value (Default is n): y

2022-03-04 12:10:38 ℹ [INFO]  Start executing the plan.
2022-03-04 12:10:38 ℹ [INFO]  Changes count: 2.
2022-03-04 12:10:38 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
2022-03-04 12:10:38 ℹ [INFO]  Processing: go-webapp-repo -> Delete ...
2022-03-04 12:10:40 ✔ [SUCCESS]  GitHub repo go-webapp-devstream-demo removed.
2022-03-04 12:10:40 ℹ [INFO]  Prepare to delete 'go-webapp-repo_github-repo-scaffolding-golang' from States.
2022-03-04 12:10:40 ✔ [SUCCESS]  Plugin go-webapp-repo delete done.
2022-03-04 12:10:40 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
2022-03-04 12:10:40 ℹ [INFO]  Processing: golang-demo-actions -> Delete ...
2022-03-04 12:10:40 ℹ [INFO]  language is go-1.17.
2022-03-04 12:10:41 ✔ [SUCCESS]  Github Actions workflow pr-builder.yml already removed.
2022-03-04 12:10:42 ✔ [SUCCESS]  Github Actions workflow main-builder.yml already removed.
2022-03-04 12:10:42 ℹ [INFO]  Prepare to delete 'golang-demo-actions_githubactions-golang' from States.
2022-03-04 12:10:42 ✔ [SUCCESS]  Plugin golang-demo-actions delete done.
2022-03-04 12:10:42 ✔ [SUCCESS]  All plugins destroyed successfully.
2022-03-04 12:10:42 ✔ [SUCCESS]  Destroy finished.
```
