# jenkins 插件

本插件使用 helm 在已有的 k8s 集群上安装 [jenkins](https://jenkins.io)。

## 使用方法

### 注意

#### 生产环境
请将配置中的 `storageClass` 修改为已存在的 storage class.

#### 测试环境
如果你想**在本地测试插件**：

1. 请将配置文件中的 `test_env` 改为 `true`。
2. 在运行 k8s 的主机上创建数据目录并修改权限，命令如下：

```bash
# 如果 k8s 和 dtm 运行在同一个主机上
mkdir -p ~/data/jenkins-volumes/
chown -R 1000:1000 ~/data/jenkins-volumes/

-------------------------

# 如果 k8s 和 dtm 运行在不同的主机上，比如 k8s 运行在 docker 中
# 1. 获取 dtm 运行的主机的用户的 home 目录
cd ~ && pwd
# 2. 复制上面的命令结果
# 3. 在 k8s 运行的主机上创建数据目录并修改权限，命令如下：
docker exec -it <container-id> bash
mkdir -p <your-dtm-home-dir>/data/jenkins-volumes/
chown -R 1000:1000 <your-dtm-home-dir>/data/jenkins-volumes/
```



### 配置

```yaml
--8<-- "jenkins.yaml"
```

当前，所有配置项均为必填。
