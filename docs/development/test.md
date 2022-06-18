# Test

### Unit Test

Run all unit tests:

```shell
go test ./...
```

### End-to-end tests

##### 1. e2e test with Github action

GitHub Actions will run e2e test automatically

##### 2. e2e test locally

```shell
bash hack/e2e/e2e-run.sh
```

This test script depends on the following environment variables:

- GITHUB_USER
- GITHUB_TOKEN
- DOCKERHUB_USERNAME
- DOCKERHUB_TOKEN

Set it before test.
