[![GitHub Release](https://img.shields.io/github/release/DocPlanner/helm-repo-updater.svg?logo=github&labelColor=262b30)](https://github.com/DocPlanner/helm-repo-updater/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/DocPlanner/helm-repo-updater)](https://goreportcard.com/report/github.com/DocPlanner/helm-repo-updater)
[![License](https://img.shields.io/github/license/DocPlanner/helm-repo-updater)](https://github.com/DocPlanner/helm-repo-updater/LICENSE)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [helm-repo-updater](#helm-repo-updater)
  - [Scope](#scope)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Examples of usage](#examples-of-usage)
  - [Running the tests](#running-the-tests)
    - [Tests requirements](#tests-requirements)
    - [Launch tests](#launch-tests)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# helm-repo-updater

## Scope

This repo aims to manage the development of `helm-repo-updater`, a CLI tool whose objective is to update one or more sets of value keys in the `values.yaml` files used by a [helm chart](https://helm.sh/docs/topics/charts/) stored in a specific git repository, using a commit with the desired change in the repository where the `values.yaml` file is stored.

## Installation

- Launch directly the binary, for that option it will be necessary go to [release page](https://github.com/DocPlanner/helm-repo-updater/releases) and download the binary needed for the architecture of the machine where it is going to run

- It is possible to use the Docker image directly, since for each [release](https://github.com/DocPlanner/helm-repo-updater/releases) the associated image will be published in the [GitHub Container Registry of the repository](ghcr.io/docplanner/helm-repo-updater).

## Usage

    Runs the helm repo updater

    Usage:
      helm-repo-updater run [flags]

    Flags:
          --app-name string                  app name
          --dry-run                          run in dry-run mode. If set to true, do not perform any changes
          --git-branch string                branch (default "develop")
          --git-commit-email string          E-Mail address to use for Git commits
          --git-commit-user string           Username to use for Git commits
          --git-dir string                   file eg. /production/charts/
          --git-file string                  file eg. values.yaml
          --git-password string              Password for github user
          --git-repo-url string              git repo url
          --helm-key-values stringToString   helm key-values sets (default [])
      -h, --help                             help for run
          --logLevel string                  set the loglevel to one of trace|debug|info|warn|error (default "info")
          --ssh-private-key string           ssh private key (only using

    Global Flags:
          --config string   config file (default is $HOME/.helm-repo-updater.yaml)

## Examples of usage

- Example run to update the `.image.tag` key to `1.1.0` in the `develop` branch of the `test-repo` repository:
  ```bash
  $ helm-repo-updater run \
    --app-name=example-app \
    --git-branch="develop" \
    --git-commit-user="test-user" \
    --git-commit-email="test-user@docplanner.com" \
    --git-file="values.yaml" \
    --helm-key-values=".image.tag=1.1.0" \
    --git-repo-url="ssh://git@localhost:2222/git-server/repos/test-repo.git" \
    --ssh-private-key="test-git-server/private_keys/helm-repo-updater-test"
  INFO[2022-03-03T16:23:10+01:00] Cloning git repository ssh://git@localhost:2222/git-server/repos/test-repo.git in temporal folder located in /var/folders/vb/v4wr_9f52ns4mmdkwp4_35cm0000gp/T/git-example-app956297090  application=example-app
  Enumerating objects: 4, done.
  Counting objects: 100% (4/4), done.
  Total 4 (delta 0), reused 0 (delta 0), pack-reused 0
  INFO[2022-03-03T16:23:11+01:00] Pulling latest changes of branch develop      application=example-app
  INFO[2022-03-03T16:23:11+01:00] Actual value for key .image.tag: 1.0.0        application=example-app
  INFO[2022-03-03T16:23:11+01:00] Setting new value for key .image.tag: 1.1.0   application=example-app
  INFO[2022-03-03T16:23:11+01:00] Adding file example-app/values.yaml to git for commit changes  application=example-app
  INFO[2022-03-03T16:23:11+01:00] It's going to commit changes with message: 🚀 automatic update of example-app
  updates key .image.tag value from '1.0.0' to '1.1.0'  application=example-app
  INFO[2022-03-03T16:23:11+01:00] It's going to push commit with hash ca9ce40520f018094a2cd7952847e7ea4bb949fe and message 🚀 automatic update of example-app
  updates key .image.tag value from '1.0.0' to '1.1.0'  application=example-app
  INFO[2022-03-03T16:23:11+01:00] Pushing changes                               application=example-app
  INFO[2022-03-03T16:23:11+01:00] Successfully pushed changes                   application=example-app
  INFO[2022-03-03T16:23:11+01:00] Successfully updated the live application spec  application=example-app
  ```

- Example run to update the `.image.tag` key to `1.1.0` in the `develop` branch of the `test-repo` repository, being `1.1.0` the value currently present in the repository for the above key:
  ```
  $ helm-repo-updater run \
    --app-name=example-app \
    --git-branch="develop" \
    --git-commit-user="test-user" \
    --git-commit-email="test-user@docplanner.com" \
    --git-file="values.yaml" \
    --helm-key-values=".image.tag=1.1.0" \
    --git-repo-url="ssh://git@localhost:2222/git-server/repos/test-repo.git" \
    --ssh-private-key="test-git-server/private_keys/helm-repo-updater-test"
  INFO[2022-03-03T16:24:17+01:00] Cloning git repository ssh://git@localhost:2222/git-server/repos/test-repo.git in temporal folder located in /var/folders/vb/v4wr_9f52ns4mmdkwp4_35cm0000gp/T/git-example-app1208822616  application=example-app
  Enumerating objects: 8, done.
  Counting objects: 100% (8/8), done.
  Compressing objects: 100% (2/2), done.
  Total 8 (delta 0), reused 0 (delta 0), pack-reused 0
  INFO[2022-03-03T16:24:17+01:00] Pulling latest changes of branch develop      application=example-app
  INFO[2022-03-03T16:24:17+01:00] Actual value for key .image.tag: 1.1.0        application=example-app
  INFO[2022-03-03T16:24:17+01:00] Setting new value for key .image.tag: 1.1.0   application=example-app
  INFO[2022-03-03T16:24:17+01:00] target for key .image.tag is the same, skipping  application=example-app
  ERRO[2022-03-03T16:24:17+01:00] Could not update application spec: nothing to update, skipping commit  application=example-app
  ERRO[2022-03-03T16:24:17+01:00] Error trying to update the example-app application: nothing to update, skipping commit  application=example-app
  ```

## Running the tests

Several tests have been created, it has been taken into account that the main functionality requires interacting with a git server, so we have implemented one with the minimum functionality through a [Docker](https://www.docker.com/) container, the files for it are present in the [test-git-server](./test-git-server/) folder.

### Tests requirements

The following software must be installed to run the tests:
- [Docker](https://docs.docker.com/get-docker/)
- [docker-compose](https://docs.docker.com/compose/install/)

### Launch tests

- Execute unit tests:
  ```
  make test
  ```

- Execute unit tests with coverage:
  ```
  make test-coverage
  ```

**It is strongly recommended that before launching the tests, the command be executed to "clean" the previous executions**:
```
make clean-test-deps
```
> The above command will recreate the container created for the git server used in the tests, so that it will start from the initial scenario expected at the beginning of the tests.
