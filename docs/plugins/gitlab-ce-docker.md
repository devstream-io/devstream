# gitlab-ce-docker Plugin

This plugin installs [GitLab](https://about.gitlab.com/) CE(Community Edition) on Docker.

## Usage

Note: 
1. the user you are using must be `root` or is in `docker` group;
2. `https` isn't supported now.

```yaml
--8<-- "gitlab-ce-docker.yaml"
```

## Some Commands May Help You

- get the password of user `root` in gitlab-ce-docker

```shell
docker exec -it gitlab grep 'Password:' /etc/gitlab/initial_root_password
```

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
