# 开发工作流

这篇文档是关于如何参与DevStream开发的流程。

## 第一步 - Fork 仓库
1. 打开项目仓库： https://github.com/devstream-io/devstream ；
2. 点击 `Fork` 按钮，从DevStream创建一个fork。

## 第二步 - Clone 仓库
1. 定义一些基础的环境变量

请根据你的实际情况来设置值。
```sh
export WORKING_PATH="~/gocode"
export USER="daniel-hutao"
export PROJECT="devstream"
export ORG="devstream-io"
```

2. Clone仓库到你本地
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

3. 确认你的远程设置是正确的

执行`git remote -v` 这个命令，你将看到如下输出：
```sh
origin	git@github.com:daniel-hutao/devstream.git (fetch)
origin	git@github.com:daniel-hutao/devstream.git (push)
upstream	https://github.com/devstream-io/devstream (fetch)
upstream	no_push (push)
```

## 第三步 - 分支代码保持同步更新

你经常需要更新你的本地代码，以便与上游保持同步。
```sh
git fetch upstream
git checkout main
git rebase upstream/main
```

## 第四步 - 编码

首先，你需要拉一个新的分支，名字根据你自己的喜好而定。

```sh
git checkout -b feat-xxx
```

然后开始编码吧。

## 第五步 - 提交&推送

```sh
git add <file>
git commit -s -m "some description here"
git push -f origin feat-xxx
```

然后你就可以在GitHub上创建一个`pr`。
