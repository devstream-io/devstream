package golang

var prPipeline = `name: Pull Requests Builder
on:
  pull_request:
    branches: [ master, main ]
permissions: write-all
jobs:
  [[- if .Build.Enable]]
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.17
    - name: Build
      run: [[- if not .Build.Command]] go build ./...[[- else]] [[.Build.Command]][[- end]]
  [[- else]]
  [[- end]]
  [[- if .Test.Enable]]
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.17
    - name: Test
      run: [[- if not .Test.Command]] go test ./...[[- else]] [[.Test.Command]][[- end]] [[- if .Test.Coverage.Enable]] -race -covermode=atomic -coverprofile=[[- if not .Test.Coverage.Output]]coverage.out[[- else]][[.Test.Coverage.Output]][[- end]] [[- end]]
    [[- if .Test.Coverage.Enable]]
    - name: Generate Coverage
      run: go tool cover -func=[[- if not .Test.Coverage.Output]]coverage.out[[- else]][[.Test.Coverage.Output]][[- end]] >> coverage.cov
    - name: comment PR
      uses: machine-learning-apps/pr-comment@master
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      with:
        path: coverage.cov
    [[- end]]
  [[- else]]
  [[- end]]
  tag:
    name: Tag
    needs: [test]
    if: ${{ github.event_name == 'push' }}
    runs-on: ubuntu-latest
    outputs:
      new_tag: ${{ steps.tag_version.outputs.new_tag }}
    steps:
      - uses: actions/checkout@v2
      - name: Bump version and push tag
        id: tag_version
        uses: mathieudutour/github-tag-action@v5.6
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          tag_prefix: ""
  [[- if .Docker.Enable]]
  image:
    name: Build Docker Image
    needs: [tag]
    if: ${{ github.event_name == 'push' }}
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Login to Docker Hub
        uses: docker/login-action@v1
        with:
          username: [[.Docker.Registry.Username]]
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      - name: Build and push
        id: docker_build
        uses: docker/build-push-action@v2
        with:
          push: true
          tags: [[.Docker.Registry.Username]]/[[- if not .Docker.Registry.Repository]][[.Repo]][[- else]][[.Docker.Registry.Repository]][[- end]]:${{needs.tag.outputs.new_tag}}
  [[- end]]
`

var mainPipeline = `name: Main Branch Builder
on:
  push:
    branches: [ master, main ]
permissions: write-all
jobs:
  [[- if .Build.Enable]]
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.17
    - name: Build
      run: [[- if not .Build.Command]] go build ./...[[- else]] [[.Build.Command]][[- end]]
  [[- else]]
  [[- end]]
  [[- if .Test.Enable]]
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.17
    - name: Test
      run: [[- if not .Test.Command]] go test ./...[[- else]] [[.Test.Command]][[- end]] [[- if .Test.Coverage.Enable]] -race -covermode=atomic -coverprofile=[[- if not .Test.Coverage.Output]]coverage.out[[- else]][[.Test.Coverage.Output]][[- end]] [[- end]]
    [[- if .Test.Coverage.Enable]]
    - name: Generate Coverage
      run: go tool cover -func=[[- if not .Test.Coverage.Output]]coverage.out[[- else]][[.Test.Coverage.Output]][[- end]] >> coverage.cov
    - name: Get comment body
      id: get-comment-body
      run: |
        cat coverage.cov
        body=$(cat coverage.cov)
        body="${body//'%'/'%25'}"
        body="${body//$'\n'/'%0A'}"
        body="${body//$'\r'/'%0D'}" 
        echo ::set-output name=body::$body
    - name: Create commit comment
      uses: peter-evans/commit-comment@v1
      with:
       body: ${{ steps.get-comment-body.outputs.body }}
    [[- end]]
  [[- else]]
  [[- end]]
  tag:
    name: Tag
    needs: [test]
    if: ${{ github.event_name == 'push' }}
    runs-on: ubuntu-latest
    outputs:
      new_tag: ${{ steps.tag_version.outputs.new_tag }}
    steps:
      - uses: actions/checkout@v2
      - name: Bump version and push tag
        id: tag_version
        uses: mathieudutour/github-tag-action@v5.6
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          tag_prefix: ""
  [[- if .Docker.Enable]]
  image:
    name: Build Docker Image
    needs: [tag]
    if: ${{ github.event_name == 'push' }}
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Login to Docker Hub
        uses: docker/login-action@v1
        with:
          username: [[.Docker.Registry.Username]]
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      - name: Build and push
        id: docker_build
        uses: docker/build-push-action@v2
        with:
          push: true
          tags: [[.Docker.Registry.Username]]/[[- if not .Docker.Registry.Repository]][[.Repo]][[- else]][[.Docker.Registry.Repository]][[- end]]:${{needs.tag.outputs.new_tag}}
  [[- end]]
`
