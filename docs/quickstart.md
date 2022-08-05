# Quick Start

[Watch the demo](./index.md) first if you prefer to see DevStream in action.

> Note: currently, we only have Linux and macOS versions of DevStream. Windows support will come later.

In this quickstart, we will do the following automatically with DevStream:

- create a GitHub repository with Golang web app scaffolding;
- set up GitHub Actions workflow for the Golang app we created, which contains test and build stages for our Go web app.

## 1 Configuration

As aforementioned, we will handle GitHub repo scaffolding and CI workflows in GitHub Actions, so, we will need the following environment variables (env vars) to be set:

- GITHUB_USER
- GITHUB_TOKEN
- DOCKERHUB_USERNAME

Run the following commands to set values for these env vars (replace values within the double quotes):

```shell
export GITHUB_USER="<YOUR_GITHUB_USER_NAME_HERE>"
export GITHUB_TOKEN="<YOUR_GITHUB_PERSONAL_ACCESS_TOKEN_HERE>"
export DOCKERHUB_USERNAME="<YOUR_DOCKER_HUB_USER_NAME_HERE>"
```

> Tip: Go to [Personal Access Token](https://github.com/settings/tokens/new) to generate a new `GITHUB_TOKEN` for `dtm`.
> 
> For "Quick Start", we only need `repo`,`workflow`,`delete_repo` scopes. However, we recommend that you check all and future plugins may require more scopes.

Then we run the following commands to update our config file with those env vars:

For **macOS** or **FreeBSD** based systems:

```shell
sed -i.bak "s@YOUR_GITHUB_USERNAME_CASE_SENSITIVE@${GITHUB_USER}@g" quickstart.yaml
sed -i.bak "s@YOUR_DOCKER_USERNAME@${DOCKERHUB_USERNAME}@g" quickstart.yaml
```

For **GNU** Linux users:

```shell
sed -i "s@YOUR_GITHUB_USERNAME_CASE_SENSITIVE@${GITHUB_USER}@g" quickstart.yaml
sed -i "s@YOUR_DOCKER_USERNAME@${DOCKERHUB_USERNAME}@g" quickstart.yaml
```

## 2 Download

In your working directory, run:

```shell
sh -c "$(curl -fsSL https://raw.githubusercontent.com/devstream-io/devstream/main/hack/quick-start/quickstart.sh)"
```

This will download the `dtm` binary and a `quickstart.yaml` config file to your working directory, and grant the binary execution permission.

> Optional: you can then move `dtm` to a place which is in your PATH. For example: `mv dtm /usr/local/bin/`.


## 3 Init

Run:

```shell
./dtm init -f quickstart.yaml
```

and you should see similar output to the following:

```
2022-06-30 11:21:48 ℹ [INFO]  Got Backend from config: local
2022-06-30 11:21:48 ℹ [INFO]  Using dir <.devstream> to store plugins.
2022-06-30 11:21:48 ℹ [INFO]  Downloading: [repo-scaffolding-darwin-arm64_0.7.0.so] ...
 15.05 MiB / 15.05 MiB [================================] 100.00% 21.17 MiB/s 0s
2022-06-30 11:21:49 ✔ [SUCCESS]  [repo-scaffolding-darwin-arm64_0.7.0.so] download succeeded.
2022-06-30 11:21:49 ℹ [INFO]  Downloading: [repo-scaffolding-darwin-arm64_0.7.0.md5] ...
 33 B / 33 B [==========================================] 100.00% 35.29 KiB/s 0s
2022-06-30 11:21:49 ✔ [SUCCESS]  [repo-scaffolding-darwin-arm64_0.7.0.md5] download succeeded.
2022-06-30 11:21:49 ℹ [INFO]  Plugin: repo-scaffolding-darwin-arm64_0.7.0.so doesn't match with .md5 and will be downloaded.
2022-06-30 11:21:49 ℹ [INFO]  Downloading: [repo-scaffolding-darwin-arm64_0.7.0.so] ...
 15.05 MiB / 15.05 MiB [================================] 100.00% 31.25 MiB/s 0s
2022-06-30 11:21:50 ✔ [SUCCESS]  [repo-scaffolding-darwin-arm64_0.7.0.so] download succeeded.
2022-06-30 11:21:50 ℹ [INFO]  Downloading: [repo-scaffolding-darwin-arm64_0.7.0.md5] ...
 33 B / 33 B [==========================================] 100.00% 43.43 KiB/s 0s
2022-06-30 11:21:50 ✔ [SUCCESS]  [repo-scaffolding-darwin-arm64_0.7.0.md5] download succeeded.
2022-06-30 11:21:50 ℹ [INFO]  Downloading: [githubactions-golang-darwin-arm64_0.7.0.so] ...
 17.49 MiB / 17.49 MiB [================================] 100.00% 31.18 MiB/s 0s
2022-06-30 11:21:51 ✔ [SUCCESS]  [githubactions-golang-darwin-arm64_0.7.0.so] download succeeded.
2022-06-30 11:21:51 ℹ [INFO]  Downloading: [githubactions-golang-darwin-arm64_0.7.0.md5] ...
 33 B / 33 B [=========================================] 100.00% 160.70 KiB/s 0s
2022-06-30 11:21:51 ✔ [SUCCESS]  [githubactions-golang-darwin-arm64_0.7.0.md5] download succeeded.
2022-06-30 11:21:51 ℹ [INFO]  Plugin: githubactions-golang-darwin-arm64_0.7.0.so doesn't match with .md5 and will be downloaded.
2022-06-30 11:21:51 ℹ [INFO]  Downloading: [githubactions-golang-darwin-arm64_0.7.0.so] ...
 17.49 MiB / 17.49 MiB [================================] 100.00% 31.78 MiB/s 0s
2022-06-30 11:21:52 ✔ [SUCCESS]  [githubactions-golang-darwin-arm64_0.7.0.so] download succeeded.
2022-06-30 11:21:52 ℹ [INFO]  Downloading: [githubactions-golang-darwin-arm64_0.7.0.md5] ...
 33 B / 33 B [==========================================] 100.00% 87.12 KiB/s 0s
2022-06-30 11:21:52 ✔ [SUCCESS]  [githubactions-golang-darwin-arm64_0.7.0.md5] download succeeded.
2022-06-30 11:21:52 ✔ [SUCCESS]  Initialize finished.
```

## 4 Apply

Run:

```shell
./dtm apply -f quickstart.yaml
```

When it prompts:

```shell
...(omitted)
Continue? [y/n]
Enter a value (Default is n):
```

input `y` and hit the enter key.

You should see similar output to the following

```
2022-06-30 11:25:47 ℹ [INFO]  Apply started.
2022-06-30 11:25:47 ℹ [INFO]  Got Backend from config: local
2022-06-30 11:25:47 ℹ [INFO]  Using dir <.devstream> to store plugins.
2022-06-30 11:25:47 ℹ [INFO]  Using local backend. State file: devstream.state.
2022-06-30 11:25:47 ℹ [INFO]  Tool (repo-scaffolding/default) found in config but doesn't exist in the state, will be created.
2022-06-30 11:25:47 ℹ [INFO]  Tool (githubactions-golang/default) found in config but doesn't exist in the state, will be created.
Continue? [y/n]
Enter a value (Default is n): y

2022-06-30 11:26:20 ℹ [INFO]  Start executing the plan.
2022-06-30 11:26:20 ℹ [INFO]  Changes count: 2.
2022-06-30 11:26:20 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
2022-06-30 11:26:20 ℹ [INFO]  Processing: (repo-scaffolding/default) -> Create ...
2022-06-30 11:26:24 ℹ [INFO]  The repo go-webapp-devstream-demo has been created.
2022-06-30 11:26:37 ✔ [SUCCESS]  Tool (repo-scaffolding/default) Create done.
2022-06-30 11:26:37 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
2022-06-30 11:26:37 ℹ [INFO]  Processing: (githubactions-golang/default) -> Create ...
2022-06-30 11:26:38 ℹ [INFO]  Creating GitHub Actions workflow pr-builder.yml ...
2022-06-30 11:26:38 ✔ [SUCCESS]  Github Actions workflow pr-builder.yml created.
2022-06-30 11:26:38 ℹ [INFO]  Creating GitHub Actions workflow main-builder.yml ...
2022-06-30 11:26:39 ✔ [SUCCESS]  Github Actions workflow main-builder.yml created.
2022-06-30 11:26:39 ✔ [SUCCESS]  Tool (githubactions-golang/default) Create done.
2022-06-30 11:26:39 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
2022-06-30 11:26:39 ✔ [SUCCESS]  All plugins applied successfully.
2022-06-30 11:26:39 ✔ [SUCCESS]  Apply finished.
```

## 5 Check the Results

Go to your GitHub repositories list and you can see the new repo `go-webapp-devstream-demo` has been created.

There is scaffolding code for a Golang web app in it, with GitHub Actions CI workflow set up properly.

The commits (made by DevStream when scaffolding the repo and creating workflows) have triggered the CI, and the workflow has finished successfully, as shown in the screenshot below:

![](./images/repo-scaffolding.png)

## 6 Clean Up

Run:

```shell
./dtm delete -f quickstart.yaml
```

input `y` the same just like you did in the previous steps, and you should see similar output:

```
2022-06-30 11:31:01 ℹ [INFO]  Delete started.
2022-06-30 11:31:01 ℹ [INFO]  Got Backend from config: local
2022-06-30 11:31:01 ℹ [INFO]  Using dir <.devstream> to store plugins.
2022-06-30 11:31:01 ℹ [INFO]  Using local backend. State file: devstream.state.
2022-06-30 11:31:01 ℹ [INFO]  Tool (githubactions-golang/default) will be deleted.
2022-06-30 11:31:01 ℹ [INFO]  Tool (repo-scaffolding/default) will be deleted.
Continue? [y/n]
Enter a value (Default is n): y

2022-06-30 11:31:03 ℹ [INFO]  Start executing the plan.
2022-06-30 11:31:03 ℹ [INFO]  Changes count: 2.
2022-06-30 11:31:03 ℹ [INFO]  -------------------- [  Processing progress: 1/2.  ] --------------------
2022-06-30 11:31:03 ℹ [INFO]  Processing: (githubactions-golang/default) -> Delete ...
2022-06-30 11:31:04 ℹ [INFO]  Deleting GitHub Actions workflow pr-builder.yml ...
2022-06-30 11:31:05 ✔ [SUCCESS]  GitHub Actions workflow pr-builder.yml removed.
2022-06-30 11:31:05 ℹ [INFO]  Deleting GitHub Actions workflow main-builder.yml ...
2022-06-30 11:31:06 ✔ [SUCCESS]  GitHub Actions workflow main-builder.yml removed.
2022-06-30 11:31:06 ℹ [INFO]  Prepare to delete 'githubactions-golang_default' from States.
2022-06-30 11:31:06 ✔ [SUCCESS]  Tool (githubactions-golang/default) delete done.
2022-06-30 11:31:06 ℹ [INFO]  -------------------- [  Processing progress: 2/2.  ] --------------------
2022-06-30 11:31:06 ℹ [INFO]  Processing: (repo-scaffolding/default) -> Delete ...
2022-06-30 11:31:06 ✔ [SUCCESS]  GitHub repo go-webapp-devstream-demo removed.
2022-06-30 11:31:06 ℹ [INFO]  Prepare to delete 'repo-scaffolding_default' from States.
2022-06-30 11:31:06 ✔ [SUCCESS]  Tool (repo-scaffolding/default) delete done.
2022-06-30 11:31:06 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
2022-06-30 11:31:06 ✔ [SUCCESS]  All plugins deleted successfully.
2022-06-30 11:31:06 ✔ [SUCCESS]  Delete finished.
```

Now if you check your GitHub repo list again, everything has been nuked by DevStream. Hooray!

> Optional: you can also remove the DevStream state file (which should be empty now) by running: `rm devstream.state`.
