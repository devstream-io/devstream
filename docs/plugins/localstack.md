# localstack plugin

This plugin installs [LocalStack](https://localstack.cloud/) in Kubernetes cluster using Helm chart.

_NOTICE: currently, this plugin is advised to use in local environment only._

## Background

[LocalStack](https://localstack.cloud/) is an AWS cloud service emulator that runs in local or CI environment.

It supports a lots of AWS services, like AWS Lambda, S3, DynamoDB, Kinesis, SQS, SNS, and so on!

And it make our life as a cloud native developer easier!

## Usage

### Config

```yaml
--8<-- "localstack.yaml"
```

### Testing

```bash
kubectl port-forward <POD> 4566:4566

curl http://localhost:4566/health
```

## Notice

If you use [KinD](https://kind.sigs.k8s.io/) in [Docker Desktop 4.3.0+](https://www.docker.com/products/docker-desktop/), you should ensure the docker engine to use Cgroups v1.

See [KinD | Known Issues](https://kind.sigs.k8s.io/docs/user/known-issues/).
