package plugin

var GithubRepoScaffoldingGolangDefaultConfig = `tools:
# name of the instance with github-repo-scaffolding-golang
- name: go-webapp-repo
  # name of the plugin
  plugin: github-repo-scaffolding-golang
  # options for the plugin
  options:
    # the repo's owner. It should be case-sensitive here; strictly use your GitHub user name; please change the value below.
    owner: YOUR_GITHUB_USERNAME
    # the repo which you'd like to create; please change the value below.
    repo: YOUR_REPO_NAME
    # the branch of the repo you'd like to hold the code
    branch: main
    # the image repo you'd like to push the container image; please change the value below.
    image_repo: YOUR_DOCKERHUB_USERNAME/YOUR_DOCKERHUB_IMAGE_REPO_NAME`
