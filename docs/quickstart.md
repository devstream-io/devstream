# Quick Start

In this quickstart, we will do the following automatically with DevStream:

- create a GitHub repository with automatically generated code for a web application written in Golang with the [gin](https://github.com/gin-gonic/gin) framework;
- set up GitHub Actions workflow for the app created in the previous step.

---

## 1 Download

In your working directory, run:

```shell
sh -c "$(curl -fsSL https://raw.githubusercontent.com/devstream-io/devstream/main/hack/install/download.sh)"
```
> Note: the command above does the following:
> 
> - find out your OS and chip architecture
> - find the latest version of the `dtm` binary
> - download the correct `dtm` according to OS/architecture
> - grant the binary execution permission.
> 
> Optional: you can then move `dtm` to a place which is in your PATH. For example: `mv dtm /usr/local/bin/`.
> 
> For more ways to install `dtm`, see [install dtm](./install.md).

---

## 2 Configuration

Run the following command to generate the template configuration file `config.yaml` for quickstart.

```shell
./dtm show config -t quickstart > config.yaml
```

Then set the following environment variables by running (replace values within the double quotes):

```shell
export GITHUB_USER="<YOUR_GITHUB_USER_NAME_HERE>"
export GITHUB_TOKEN="<YOUR_GITHUB_PERSONAL_ACCESS_TOKEN_HERE>"
export DOCKERHUB_USERNAME="<YOUR_DOCKER_HUB_USER_NAME_HERE>"
```

> Tip: go to [Personal Access Token](https://github.com/settings/tokens/new) to generate a new `GITHUB_TOKEN` for `dtm`.
> 
> For "Quick Start", we only need `repo`,`workflow`,`delete_repo` permissions.

Then we run the following commands to update our config file with those env vars:

Then, if you are using **macOS** or **FreeBSD** based systems, run:

```shell
sed -i.bak "s@YOUR_GITHUB_USERNAME_CASE_SENSITIVE@${GITHUB_USER}@g" config.yaml
sed -i.bak "s@YOUR_DOCKER_USERNAME@${DOCKERHUB_USERNAME}@g" config.yaml
```

> For **GNU** Linux users, run:
> 
> ```shell
> sed -i "s@YOUR_GITHUB_USERNAME_CASE_SENSITIVE@${GITHUB_USER}@g" config.yaml
> sed -i "s@YOUR_DOCKER_USERNAME@${DOCKERHUB_USERNAME}@g" config.yaml
> ```

---

## 3 Init

Run:

```shell
./dtm init
```

> You should see some output similar to the following:
```
2022-12-02 16:11:55 ℹ [INFO]  Using dir </Users/tiexin/.devstream/plugins> to store plugins.
2022-12-02 16:11:55 ℹ [INFO]  -------------------- [  repo-scaffolding-darwin-arm64_0.10.1  ] --------------------
2022-12-02 16:11:57 ℹ [INFO]  Downloading: [repo-scaffolding-darwin-arm64_0.10.1.so] ...
 87.82 MiB / 87.82 MiB [================================] 100.00% 12.30 MiB/s 7s
2022-12-02 16:12:04 ✔ [SUCCESS]  [repo-scaffolding-darwin-arm64_0.10.1.so] download succeeded.
2022-12-02 16:12:04 ℹ [INFO]  Downloading: [repo-scaffolding-darwin-arm64_0.10.1.md5] ...
 33 B / 33 B [==========================================] 100.00% 50.98 KiB/s 0s
2022-12-02 16:12:04 ✔ [SUCCESS]  [repo-scaffolding-darwin-arm64_0.10.1.md5] download succeeded.
2022-12-02 16:12:04 ℹ [INFO]  Initialize [repo-scaffolding-darwin-arm64_0.10.1] finished.
2022-12-02 16:12:04 ℹ [INFO]  -------------------- [  repo-scaffolding-darwin-arm64_0.10.1  ] --------------------
2022-12-02 16:12:04 ℹ [INFO]  -------------------- [  githubactions-golang-darwin-arm64_0.10.1  ] --------------------
2022-12-02 16:12:05 ℹ [INFO]  Downloading: [githubactions-golang-darwin-arm64_0.10.1.so] ...
 86.44 MiB / 86.44 MiB [================================] 100.00% 15.12 MiB/s 5s
2022-12-02 16:12:10 ✔ [SUCCESS]  [githubactions-golang-darwin-arm64_0.10.1.so] download succeeded.
2022-12-02 16:12:10 ℹ [INFO]  Downloading: [githubactions-golang-darwin-arm64_0.10.1.md5] ...
 33 B / 33 B [==========================================] 100.00% 71.24 KiB/s 0s
2022-12-02 16:12:10 ✔ [SUCCESS]  [githubactions-golang-darwin-arm64_0.10.1.md5] download succeeded.
2022-12-02 16:12:11 ℹ [INFO]  Initialize [githubactions-golang-darwin-arm64_0.10.1] finished.
2022-12-02 16:12:11 ℹ [INFO]  -------------------- [  githubactions-golang-darwin-arm64_0.10.1  ] --------------------
2022-12-02 16:12:11 ✔ [SUCCESS]  Initialize finished.
```

---

## 4 Apply

Run:

```shell
./dtm apply -y
```

> You should see similar output to the following

```
2022-12-02 16:18:00 ℹ [INFO]  Apply started.
2022-12-02 16:18:00 ℹ [INFO]  Using local backend. State file: devstream.state.
2022-12-02 16:18:00 ℹ [INFO]  Tool (repo-scaffolding/golang-github) found in config but doesn't exist in the state, will be created.
2022-12-02 16:18:00 ℹ [INFO]  Tool (githubactions-golang/default) found in config but doesn't exist in the state, will be created.
2022-12-02 16:18:00 ℹ [INFO]  Start executing the plan.
2022-12-02 16:18:00 ℹ [INFO]  Changes count: 2.
2022-12-02 16:18:00 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
2022-12-02 16:18:00 ℹ [INFO]  Processing: (repo-scaffolding/golang-github) -> Create ...
2022-12-02 16:18:00 ℹ [INFO]  github start to download repoTemplate...
2022-12-02 16:18:04 ✔ [SUCCESS]  The repo go-webapp-devstream-demo has been created.
2022-12-02 16:18:12 ✔ [SUCCESS]  Tool (repo-scaffolding/golang-github) Create done.
2022-12-02 16:18:12 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
2022-12-02 16:18:12 ℹ [INFO]  Processing: (githubactions-golang/default) -> Create ...
2022-12-02 16:18:13 ℹ [INFO]  Creating GitHub Actions workflow pr-builder.yml ...
2022-12-02 16:18:14 ✔ [SUCCESS]  Github Actions workflow pr-builder.yml created.
2022-12-02 16:18:14 ℹ [INFO]  Creating GitHub Actions workflow main-builder.yml ...
2022-12-02 16:18:15 ✔ [SUCCESS]  Github Actions workflow main-builder.yml created.
2022-12-02 16:18:15 ✔ [SUCCESS]  Tool (githubactions-golang/default) Create done.
2022-12-02 16:18:15 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
2022-12-02 16:18:15 ✔ [SUCCESS]  All plugins applied successfully.
2022-12-02 16:18:15 ✔ [SUCCESS]  Apply finished.
```

---

## 5 Check the Results

Go to your GitHub repositories list and you can see the new repo `go-webapp-devstream-demo` has been created.

There is scaffolding code for a Golang web app in it, with GitHub Actions CI workflow set up properly.

The commits (made by DevStream when scaffolding the repo and creating workflows) have triggered the CI, and the workflow has finished successfully, as shown in the screenshot below:

![](./images/repo-scaffolding.png)

---

## 6 Clean Up

Run:

```shell
./dtm delete
```

Input `y` then press enter to continue, and you should see similar output:

```
2022-12-02 16:19:07 ℹ [INFO]  Delete started.
2022-12-02 16:19:07 ℹ [INFO]  Using local backend. State file: devstream.state.
2022-12-02 16:19:07 ℹ [INFO]  Tool (githubactions-golang/default) will be deleted.
2022-12-02 16:19:07 ℹ [INFO]  Tool (repo-scaffolding/golang-github) will be deleted.
Continue? [y/n]
Enter a value (Default is n): y

2022-12-02 16:19:08 ℹ [INFO]  Start executing the plan.
2022-12-02 16:19:08 ℹ [INFO]  Changes count: 2.
2022-12-02 16:19:08 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
2022-12-02 16:19:08 ℹ [INFO]  Processing: (githubactions-golang/default) -> Delete ...
2022-12-02 16:19:09 ℹ [INFO]  Deleting GitHub Actions workflow pr-builder.yml ...
2022-12-02 16:19:09 ✔ [SUCCESS]  GitHub Actions workflow pr-builder.yml removed.
2022-12-02 16:19:10 ℹ [INFO]  Deleting GitHub Actions workflow main-builder.yml ...
2022-12-02 16:19:10 ✔ [SUCCESS]  GitHub Actions workflow main-builder.yml removed.
2022-12-02 16:19:10 ℹ [INFO]  Prepare to delete 'githubactions-golang_default' from States.
2022-12-02 16:19:10 ✔ [SUCCESS]  Tool (githubactions-golang/default) delete done.
2022-12-02 16:19:10 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
2022-12-02 16:19:10 ℹ [INFO]  Processing: (repo-scaffolding/golang-github) -> Delete ...
2022-12-02 16:19:11 ✔ [SUCCESS]  GitHub repo go-webapp-devstream-demo removed.
2022-12-02 16:19:11 ℹ [INFO]  Prepare to delete 'repo-scaffolding_golang-github' from States.
2022-12-02 16:19:11 ✔ [SUCCESS]  Tool (repo-scaffolding/golang-github) delete done.
2022-12-02 16:19:11 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
2022-12-02 16:19:11 ✔ [SUCCESS]  All plugins deleted successfully.
2022-12-02 16:19:11 ✔ [SUCCESS]  Delete finished.
```

Now if you check your GitHub repo list again, everything has been nuked by DevStream. Hooray!

You can also remove the DevStream state file (which should be empty now) by running: `rm devstream.state`.
