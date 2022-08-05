# gitlab-ce-docker 插件

这个插件用于以 Docker 的方式安装 [GitLab](https://about.gitlab.com/) CE（社区版）。

## 用法

注意: 
1. 你使用的用户必须是 `root` 或者在 `docker` 用户组里；
2. 目前暂不支持 `https` 方式访问 GitLab。

```yaml
--8<-- "gitlab-ce-docker.yaml"
```

## 一些可能有用的命令

- 查看 gitlab 的 root 用户的密码:

```shell
sudo docker exec -it gitlab grep 'Password:' /etc/gitlab/initial_root_password
```

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
