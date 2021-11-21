package githubactions

var prBuilder = `
name: PR Builder
on:
  pull_request:
    branches: [ master ]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.17
    - name: Test
      run: go test -v ./...
    - name: Build
      run: go build -v ./...
`

var masterBuilder = `
name: Master Builder
on:
  push:
    branches: [ master ]
jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.17
      - name: Test
        run: go test -v ./...
      - name: Build
        run: go build -v ./...
  tag:
    name: Tag
    needs: [test]
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
  build:
    name: Build Docker Image
    needs: [tag]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Login to Docker Hub
        uses: docker/login-action@v1
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      - name: Build and push
        id: docker_build
        uses: docker/build-push-action@v2
        with:
          push: true
          tags: ${{ secrets.DOCKERHUB_USERNAME }}/go-hello-http:${{needs.tag.outputs.new_tag}}
`
