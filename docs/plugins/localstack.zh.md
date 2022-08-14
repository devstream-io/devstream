# localstack 插件

这个插件使用 Helm 在 K8S 集群上安装 [LocalStack](LocalStack)。

注意：目前仅建议在本地环境中使用本插件。

## 背景

[LocalStack](https://localstack.cloud/)是一个 AWS 云服务模拟器，可以在本地或 CI 环境中运行。

它支持大量的 AWS 服务，如 AWS Lambda, S3, DynamoDB, Kinesis, SQS, SNS，……!

它能使云原生开发者的生活更轻松!

## 用例

### 配置

```yaml
--8<-- "localstack.yaml"
```

### 测试

```bash
kubectl port-forward <POD> 4566:4566

curl http://localhost:4566/health
```

## 注意

如果你在 [Docker Desktop 4.3.0+](https://www.docker.com/products/docker-desktop/) 中使用 [KinD](https://kind.sigs.k8s.io/)，需要确保 Docker 引擎使用 Cgroups v1.

详见 [KinD | Known Issues](https://kind.sigs.k8s.io/docs/user/known-issues/).
