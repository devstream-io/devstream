# gitlab-ce-docker Plugin

This plugin installs [GitLab](https://about.gitlab.com/) CE(Community Edition) on Docker.

_NOTICE: currently, this plugin support Linux only._

## 背景知识

GitLab officially provides an image [gitlab-ce](https://registry.hub.docker.com/r/gitlab/gitlab-ce). We can use this image to start a container:

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

The variable $GITLAB_HOME here pointing to the directory where the configuration, logs, and data files will reside.
We could set this variable by the `export` command:

```shell
export GITLAB_HOME=/srv/gitlab
```

The GitLab container uses host mounted volumes to store persistent data:

| Local location        |Container location |                    Usage                   |
| --------------------- | ----------------- | ------------------------------------------ |
| `$GITLAB_HOME/data`   | `/var/opt/gitlab` | For storing application data               |
| `$GITLAB_HOME/logs`   | `/var/log/gitlab` | For storing logs                           |
| `$GITLAB_HOME/config` | `/etc/gitlab`     | For storing the GitLab configuration files |

So, we can customize some configurations as follows:

1. hostname
2. host port
3. persistent data path
4. docker image tag

## Configuration

Note: 
1. the user you are using must be `root` or is in `docker` group;
2. `https` isn't supported now.

```yaml
--8<-- "gitlab-ce-docker.yaml"
```

## Some Commands May Help You

- clone code

```shell
export hostname=YOUR_HOSTNAME
export username=YOUR_USERNAME
export project=YOUR_PROJECT_NAME
```

1. ssh

```shell
# port is 22
git clone git@${hostname}/${username}/${project}.git
# port is not 22, 2022 as a sample
git clone ssh://git@${hostname}:2022/${username}/${project}.git
```

2. http

```shell
# port is 80
git clone http://${hostname}/${username}/${project}.git
# port is not 80, 8080 as a sample
git clone http://${hostname}:8080/${username}/${project}.git
```
