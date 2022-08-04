package python

var lintPipeline = `name: Lint

on: push

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: autopep8
        uses: peter-evans/autopep8@v1
        with:
          args: --recursive --in-place --aggressive --aggressive .
      - name: Create Pull Request
        uses: peter-evans/create-pull-request@v3
        with:
          commit-message: autopep8 action fixes
          title: Fixes by autopep8 action
          body: This is an auto-generated PR with fixes by autopep8.
          labels: autopep8, automated pr
          branch: autopep8-patches
`

var testPipeline = `name: PR Builder

on:
  pull_request:
    branches: [ master, main ]
  push:
    branches: [ master, main ]

permissions: write-all

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Install dependencies
        run: |
          python -m pip install --upgrade pip
          pip install -r requirements.txt
      - name: Test
        run: |
          python3 -m unittest
`

var dockerPipeline = `name: Docker Builder

on:
  push:
    branches: [ main ]

jobs:
  tag:
    name: Tag
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
`
