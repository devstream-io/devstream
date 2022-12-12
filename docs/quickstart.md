# Quick Start

In this quickstart, we will do the following automatically with DevStream:

- create a GitHub repository with automatically generated code for a web application written in Golang with the [gin](https://github.com/gin-gonic/gin) framework;
- set up GitHub Actions workflow for the app created in the previous step.

---

## 1 Download

In your working directory, run:

```shell
sh -c "$(curl -fsSL https://download.devstream.io/download.sh)"
```

!!! note "Note"
    the command above does the following:
 
    - find out your OS and chip architecture
    - find the latest version of the `dtm` binary
    - download the correct `dtm` according to OS/architecture
    - grant the binary execution permission.

!!! quote "Optional"
    You can then move `dtm` to a place which is in your PATH. For example: `mv dtm /usr/local/bin/`.
    
    For more ways to install `dtm`, see [install dtm](./install.md).

---

## 2 Configuration

Run the following command to generate the template configuration file `config.yaml` for quickstart.

```shell
./dtm show config -t quickstart > config.yaml
```

Then set the following environment variables by running (replace values within the double quotes):

```shell
export GITHUB_TOKEN="<YOUR_GITHUB_PERSONAL_ACCESS_TOKEN_HERE>"
export IMAGE_REPO_PASSWORD="<YOUR_DOCKER_HUB_USER_NAME_HERE>"
```

!!! tip "Tip"
    Go to [Personal Access Token](https://github.com/settings/tokens/new) to generate a new `GITHUB_TOKEN` for `dtm`.
    
    For "Quick Start", we only need `repo`,`workflow`,`delete_repo` permissions.

Then we run the following commands to update our config file with those env vars:

===  "**macOS** or **FreeBSD** based systems"

    ```shell
    sed -i.bak "s@YOUR_GITHUB_USERNAME_CASE_SENSITIVE@${GITHUB_USER}@g" config.yaml
    sed -i.bak "s@YOUR_DOCKER_USERNAME@${DOCKERHUB_USERNAME}@g" config.yaml
    ```

=== "**GNU** Linux users"
 
    ```shell
    sed -i "s@YOUR_GITHUB_USERNAME_CASE_SENSITIVE@${GITHUB_USER}@g" config.yaml
    sed -i "s@YOUR_DOCKER_USERNAME@${DOCKERHUB_USERNAME}@g" config.yaml
    ```

---

## 3 Init

Run:

```shell
./dtm init -f config.yaml
```

!!! success "You should see some output similar to the following"
    ```text title=""
    2022-12-12 11:40:34 ℹ [INFO]  Using dir </Users/stein/.devstream/plugins> to store plugins.
    2022-12-12 11:40:34 ℹ [INFO]  -------------------- [  repo-scaffolding-darwin-arm64_0.10.2  ] --------------------
    2022-12-12 11:40:35 ℹ [INFO]  Downloading: [repo-scaffolding-darwin-arm64_0.10.2.so] ...
     87.75 MiB / 87.75 MiB [================================] 100.00% 7.16 MiB/s 12s
    2022-12-12 11:40:47 ✔ [SUCCESS]  [repo-scaffolding-darwin-arm64_0.10.2.so] download succeeded.
    2022-12-12 11:40:48 ℹ [INFO]  Downloading: [repo-scaffolding-darwin-arm64_0.10.2.md5] ...
     33 B / 33 B [=========================================] 100.00% 115.84 KiB/s 0s
    2022-12-12 11:40:48 ✔ [SUCCESS]  [repo-scaffolding-darwin-arm64_0.10.2.md5] download succeeded.
    2022-12-12 11:40:48 ℹ [INFO]  Initialize [repo-scaffolding-darwin-arm64_0.10.2] finished.
    2022-12-12 11:40:48 ℹ [INFO]  -------------------- [  repo-scaffolding-darwin-arm64_0.10.2  ] --------------------

    2022-12-12 11:40:48 ℹ [INFO]  -------------------- [  github-actions-darwin-arm64_0.10.2  ] --------------------
    2022-12-12 11:40:48 ℹ [INFO]  Downloading: [github-actions-darwin-arm64_0.10.2.so] ...
     90.27 MiB / 90.27 MiB [================================] 100.00% 10.88 MiB/s 8s
    2022-12-12 11:40:57 ✔ [SUCCESS]  [github-actions-darwin-arm64_0.10.2.so] download succeeded.
    2022-12-12 11:40:57 ℹ [INFO]  Downloading: [github-actions-darwin-arm64_0.10.2.md5] ...
     33 B / 33 B [=========================================] 100.00% 145.46 KiB/s 0s
    2022-12-12 11:40:57 ✔ [SUCCESS]  [github-actions-darwin-arm64_0.10.2.md5] download succeeded.
    2022-12-12 11:40:57 ℹ [INFO]  Initialize [github-actions-darwin-arm64_0.10.2] finished.
    2022-12-12 11:40:57 ℹ [INFO]  -------------------- [  github-actions-darwin-arm64_0.10.2  ] --------------------
    2022-12-12 11:40:57 ✔ [SUCCESS]  Initialize finished.
    ```

---

## 4 Apply

Run:

```shell
./dtm apply -f config.yaml -y
```

!!! success "You should see similar output to the following"

    ```text title=""
    2022-12-12 11:44:39 ℹ [INFO]  Apply started.
    2022-12-12 11:44:39 ℹ [INFO]  Using local backend. State file: devstream.state.
    2022-12-12 11:44:39 ℹ [INFO]  Tool (repo-scaffolding/golang-github) found in config but doesn't exist in the state, will be created.
    2022-12-12 11:44:39 ℹ [INFO]  Tool (github-actions/default) found in config but doesn't exist in the state, will be created.
    2022-12-12 11:44:39 ℹ [INFO]  Start executing the plan.
    2022-12-12 11:44:39 ℹ [INFO]  Changes count: 2.
    2022-12-12 11:44:39 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
    2022-12-12 11:44:39 ℹ [INFO]  Processing: (repo-scaffolding/golang-github) -> Create ...
    2022-12-12 11:44:39 ℹ [INFO]  github start to download repoTemplate...2022-12-12 11:44:42 ✔ [SUCCESS]  The repo go-webapp-devstream-demo has been created.
    2022-12-12 11:44:49 ✔ [SUCCESS]  Tool (repo-scaffolding/golang-github) Create done.
    2022-12-12 11:44:49 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
    2022-12-12 11:44:49 ℹ [INFO]  Processing: (github-actions/default) -> Create ...
    2022-12-12 11:44:57 ✔ [SUCCESS]  Tool (github-actions/default) Create done.
    2022-12-12 11:44:57 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
    2022-12-12 11:44:57 ✔ [SUCCESS]  All plugins applied successfully.
    2022-12-12 11:44:57 ✔ [SUCCESS]  Apply finished.
    ```

---

## 5 Check the Results

Go to your GitHub repositories list and you can see the new repo `go-webapp-devstream-demo` has been created.

There is scaffolding code for a Golang web app in it, with GitHub Actions CI workflow set up properly.

The commits (made by DevStream when scaffolding the repo and creating workflows) have triggered the CI, and the workflow has finished successfully, as shown in the screenshot below:

![](./images/quickstart.png)

---

## 6 Clean Up

Run:

```shell
./dtm delete -f config.yaml
```

Input `y` then press enter to continue, and you should see similar output:

!!! success "Output"

    ```text title=""
    2022-12-12 12:29:00 ℹ [INFO]  Delete started.
    2022-12-12 12:29:00 ℹ [INFO]  Using local backend. State file: devstream.state.
    2022-12-12 12:29:00 ℹ [INFO]  Tool (github-actions/default) will be deleted.
    2022-12-12 12:29:00 ℹ [INFO]  Tool (repo-scaffolding/golang-github) will be deleted.
    Continue? [y/n]
    Enter a value (Default is n): y
    2022-12-12 12:29:00 ℹ [INFO]  Start executing the plan.
    2022-12-12 12:29:00 ℹ [INFO]  Changes count: 2.
    2022-12-12 12:29:00 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
    2022-12-12 12:29:00 ℹ [INFO]  Processing: (github-actions/default) -> Delete ...
    2022-12-12 12:29:02 ℹ [INFO]  Prepare to delete 'github-actions_default' from States.
    2022-12-12 12:29:02 ✔ [SUCCESS]  Tool (github-actions/default) delete done.
    2022-12-12 12:29:02 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
    2022-12-12 12:29:02 ℹ [INFO]  Processing: (repo-scaffolding/golang-github) -> Delete ...
    2022-12-12 12:29:03 ✔ [SUCCESS]  GitHub repo go-webapp-devstream-demo removed.
    2022-12-12 12:29:03 ℹ [INFO]  Prepare to delete 'repo-scaffolding_golang-github' from States.
    2022-12-12 12:29:03 ✔ [SUCCESS]  Tool (repo-scaffolding/golang-github) delete done.
    2022-12-12 12:29:03 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
    2022-12-12 12:29:03 ✔ [SUCCESS]  All plugins deleted successfully.
    2022-12-12 12:29:03 ✔ [SUCCESS]  Delete finished.
    ```

Now if you check your GitHub repo list again, everything has been nuked by DevStream. Hooray!

You can also remove the DevStream state file (which should be empty now) by running: `rm devstream.state`.
