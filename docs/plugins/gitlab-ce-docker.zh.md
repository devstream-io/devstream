# gitlab-ce-docker 插件

这个插件用来在本机已存在的 Docker 上安装 [Gitlab-CE](https://about.gitlab.com/)
## 使用

注意: 
1. 执行本插件的用户，必须在 `docker` 用户组内，或者是 `root`
2. 目前暂不支持 `https` 访问 gitlab

```yaml

--8<-- "gitlab-ce-docker.yaml"

```

## 可能会用到的命令

查看 gitlab 的 root 用户的密码:
```shell
sudo docker exec -it gitlab grep 'Password:' /etc/gitlab/initial_root_password
```

克隆项目：
```shell
#ssh
# 22 port
git clone git@hostname/.../xxx.git
# if not 22 port
git clone ssh://git@hostname:port/.../xxx.git

# http
# 80 port
git clone http://hostname/.../xxx.git
# if not 80 port
git clone http://hostname:port/.../xxx.git
```
