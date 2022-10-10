# gitlab-ce-docker 插件

这个插件用于以 Docker 的方式安装 [GitLab](https://about.gitlab.com/) CE（社区版）。

_注意：目前本插件仅支持 Linux。_

## 背景知识

GitLab 官方提供了 [gitlab-ce](https://registry.hub.docker.com/r/gitlab/gitlab-ce) 镜像，通过这个镜像我们可以实现类似这样的命令来启动一个 GitLab 容器：

```shell
docker run --detach \
  --hostname gitlab.example.com \
  --publish 443:443 --publish 80:80 --publish 22:22 \
  --name gitlab \
  --restart always \
  --volume $GITLAB_HOME/config:/etc/gitlab \
  --volume $GITLAB_HOME/logs:/var/log/gitlab \
  --volume $GITLAB_HOME/data:/var/opt/gitlab \
  --shm-size 256m \
  gitlab/gitlab-ce:rc
```

其中 $GITLAB_HOME 表示的是本地存储卷路径，比如我们可以通过 export 命令来设置这个变量：

```shell
export GITLAB_HOME=/srv/gitlab
```

在上述命令中，我们可以看到这个容器使用了3个存储卷，含义分别如下：

| 本地路径               | 容器内路径          | 用途               |
| --------------------- | ----------------- | ----------------- |
| `$GITLAB_HOME/data`   | `/var/opt/gitlab` | 保存应用数据        |
| `$GITLAB_HOME/logs`   | `/var/log/gitlab` | 保存日志            |
| `$GITLAB_HOME/config` | `/etc/gitlab`     | 保存 GitLab 配置文件 |

在此基础上，我们可以自定义如下一些配置：

1. hostname
2. 本机端口
3. 存储卷路径
4. 镜像版本

## 配置

注意: 
1. 你使用的用户必须是 `root` 或者在 `docker` 用户组里；
2. 目前暂不支持 `https` 方式访问 GitLab。

下面的配置文件展示的是"tool file"的内容。

关于更多关于DevStream的主配置、tool file、var file的信息，请阅读[核心概念概览](../core-concepts/core-concepts.zh.md)和[DevStream配置](../core-concepts/config.zh.md).

```yaml
--8<-- "gitlab-ce-docker.yaml"
```

## 一些可能有用的命令

- 克隆项目

```shell
export hostname=YOUR_HOSTNAME
export username=YOUR_USERNAME
export project=YOUR_PROJECT_NAME
```

1. ssh 方式

```shell
# port is 22
git clone git@${hostname}/${username}/${project}.git
# port is not 22, 2022 as a sample
git clone ssh://git@${hostname}:2022/${username}/${project}.git
```

2. http 方式

```shell
# port is 80
git clone http://${hostname}/${username}/${project}.git
# port is not 80, 8080 as a sample
git clone http://${hostname}:8080/${username}/${project}.git
```
