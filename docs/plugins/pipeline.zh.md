# Pipeline 配置项

## imageRepo

若设置了该配置项，则流水线中会基于项目打包出一个 Docker 项目镜像并推送到配置的镜像仓库中。

_使用该配置项前请设置以下环境变量：_

- IMAGE_REPO_PASSWORD：改变量为镜像仓库的登陆密码或者 token，用于 Devstream 上传镜像时的权限认证。

### 具体配置字段

| 字段名 | 描述                                        |
| ------ | ------------------------------------------- |
| user   | 镜像仓库的用户名                            |
| url    | 镜像仓库的地址，默认为 Dockerhub 的官方地址 |

## dingTalk

若设置了该配置项，则会在流水线失败或者成功后往指定的钉钉群发送通知。

### 具体配置字段

| 字段名        | 描述                                                         |
| ------------- | ------------------------------------------------------------ |
| name          | 钉钉机器人的名称                                             |
| webhook       | 在钉钉群中创建机器人后获取到的回调地址                       |
| securityType  | 钉钉发送通知的认证方式，目前支持 SECRET 和 KEY               |
| securityValue | 若 securityType 配置为 SECRET，该字段表示对应的 SECRET VALUE |

## sonarqube

若设置了该配置项，则会在流水线运行测试的同时使用 [sonarqube](https://www.sonarqube.org/) 扫描代码。

### 具体配置字段

| 字段名 | 描述                                                                                                                    |
| ------ | ----------------------------------------------------------------------------------------------------------------------- |
| name   | sonarqube 的配置名称                                                                                                    |
| token  | sonarqube 的 token，获取方式可以查阅该[文档](https://sonarqube.inria.fr/sonarqube/documentation/user-guide/user-token/) |
| url    | soanrqube 的 url 地址                                                                                                   |

## language

若设置了该配置项，则会在流水线模版渲染时配置语言和框架的默认值。如设置字段 `name` 为 `golang`，则流水线中的测试镜像会使用 `golang` 镜像，命令会使用 `go test ./...`。

### 具体配置字段

| 字段名    | 描述             |
| --------- | ---------------- |
| name      | 编程语言的名称   |
| framework | 所使用的框架     |
| version   | 编程语言使用版本 |
