# Development Workflow

This document shows the workflow of how to develop DevStream.

## Step 1 - Fork the repo

1. Visit the DevStream repo: [https://github.com/devstream-io/devstream](https://github.com/devstream-io/devstream);
2. Click the `Fork` button to create a fork of the DevStream.

## Step 2 - Clone the repo

1. Define some basic environment variables

Please set the appropriate values according to your actual environment.

```sh
export WORKING_PATH="~/gocode"
export USER="daniel-hutao"
export PROJECT="devstream"
export ORG="devstream-io"
```

2. Create your clone locally

```sh
mkdir -p ${WORKING_PATH}
cd ${WORKING_PATH}
# You can also use the url: git@github.com:${USER}/${PROJECT}.git
# if your ssh configuration is proper
git clone https://github.com/${USER}/${PROJECT}.git
cd ${PROJECT}

git remote add upstream https://github.com/${ORG}/${PROJECT}.git
# Never push to upstream locally
git remote set-url --push upstream no_push
```

3. Confirm the remotes you've set is make sense

Execute `git remote -v` and you'll see output like below:

```sh
origin	git@github.com:daniel-hutao/devstream.git (fetch)
origin	git@github.com:daniel-hutao/devstream.git (push)
upstream	https://github.com/devstream-io/devstream (fetch)
upstream	no_push (push)
```

## Step 3 - Keep your branch in sync

You will often need to update your local code to keep in sync with upstream

```sh
git fetch upstream
git checkout main
git rebase upstream/main
```

## Step 4 - Coding

First, you need to pull a new branch, the name is according to your own taste.

```sh
git checkout -b feat-xxx
```

Then start coding.

## Step 5 - Commit & Push

```sh
git add <file>
git commit -s -m "some description here"
git push -f origin feat-xxx
```

Then you can create a `pr` on GitHub.
