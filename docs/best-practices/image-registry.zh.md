# 部署私有镜像仓库

如果在离线环境中使用 DevStream，你需要准备一个私有镜像仓库用于保存相关容器镜像。
如果你的环境里还没有一个可用的镜像仓库，那么可以参考本文步骤快速部署一个简易的容器镜像仓库用于支持离线使用 DevStream。

*提示：本文基于 CentOS 7 编写。*

## 一、准备证书

下文域名均以 `registry.devstream.io` 为例，你在实际执行的时候需要按需修改。

```shell
cd ~ && mkdir certs
openssl genrsa -out certs/ca.key 2048
openssl req -new -x509 -days 3650 -key certs/ca.key -subj "/C=CN/ST=GD/L=SZ/O=DevStream, Inc./CN=DevStream Root CA" -out certs/ca.crt
openssl req -newkey rsa:2048 -nodes -keyout certs/domain.key -subj "/C=CN/ST=GD/L=SZ/O=DevStream, Inc./CN=*.devstream.io" -out certs/domain.csr
openssl x509 -req -extfile <(printf "subjectAltName=DNS:devstream.io,DNS:registry.devstream.io") \
  -days 3650 -in certs/domain.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/domain.crt
```

## 二、启动 Docker Registry

```shell
docker run -d \
  --restart=always \
  --name registry \
  -v $(pwd)/certs:/certs \
  -v $(pwd)/registry:/var/lib/registry \
  -e REGISTRY_HTTP_ADDR=0.0.0.0:443 \
  -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt \
  -e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
  -p 443:443 \
  registry:2
```

# 三、配置 Docker Registry

1、在 `/etc/hosts` 中配置本机 IP 到自定义域名（如：registry.devstream.io）之间的映射，如：

```shell
# docker registry
192.168.39.100 registry.devstream.io
```

2、配置 Docker 信任刚才生成的证书（域名以 registry.devstream.io 为例）

```shell
sudo mkdir -p /etc/docker/certs.d/registry.devstream.io
sudo cp ~/certs/ca.crt /etc/docker/certs.d/registry.devstream.io/ca.crt
```

3. 验证 Docker Registry 可用

```shell
docker pull busybox
docker tag busybox registry.devstream.io/busybox
docker push registry.devstream.io/busybox
```
