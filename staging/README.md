# Staging Repos

This directory contains the staging repositories for the Devstream project.

Each sub folder under `staging/` be synchronized to the repository of the same name under devstream-io.

## How to Create a New Staging Repo

Person who wants to create a new staging repo:

1. Create a new directory under `staging/` with the name of the repo you want to create.
2. Update [`.github/sync-staging-repo.yml`](../.github/sync-staging-repo.yml).
3. Pull Request to the `main` branch of repo `devstream-io/devstream`.

Reviewers:

1. Review the PR, make sure everything is correct, and do not merge it immediately.
2. Create the repo under the `devstream-io` organization and **create a branch for the repo**. If you don't have the permission to create a repo, please ask for the members of the `devstream-io`.
3. Merge the PR.
