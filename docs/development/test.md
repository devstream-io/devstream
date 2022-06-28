# Test

## Unit Test

Run all unit tests:

```shell
# at the root of the repo
go test ./...
```

_Note: not all tests are strictly "unit" tests at the moment because some tests will actually rely on internet. Help wanted here :)_

## E2E(End-to-End) Test

### GitHub Actions

Our e2e test will run in GitHub Actions automatically when there is a change on the main branch.

The definition of the GitHub Action is [here](https://github.com/devstream-io/devstream/blob/main/.github/workflows/e2e-test.yml), and the configuration files used in e2e tests are [here](https://github.com/devstream-io/devstream/tree/main/test/e2e/yaml).

### Run E2E Test Locally

We have a simple e2e test that tests the following plugins:

- `github-repo-scaffolding-golang`
- `githubactions-golang`
- `argocd`
- `argocdapp`

The template for the config file is located [here](https://github.com/devstream-io/devstream/blob/main/test/e2e/yaml/e2e-test-local.yaml).

To run the e2e test locally, you will need docker up and running first.

Then, set the following environment variables:

- GITHUB_USER
- GITHUB_TOKEN
- DOCKERHUB_USERNAME
- DOCKERHUB_TOKEN

And execute the following command:

```shell
bash hack/e2e/e2e-run.sh
```

This test script will download kind/kubectl, start a K8s cluster as a docker container using kind, then execute dtm commands, check result, and clean up the environment.
