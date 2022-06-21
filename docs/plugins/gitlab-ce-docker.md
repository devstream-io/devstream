# gitlab-ce-docker plugin

This plugin installs [Gitlab-CE](https://about.gitlab.com/) in an existing docker.
## Usage

Note: 
1. the user must be `root` or in `docker` group.
2. https not support now(todo).

```yaml

--8<-- "gitlab-ce-docker.yaml"

```

## Next
Here are some commands that may help you:

get password of user root in gitlab-ce-docker
```shell
sudo docker exec -it gitlab grep 'Password:' /etc/gitlab/initial_root_password
```

git clone:
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
