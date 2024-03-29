[![Tests][tests-badge]][tests-link]
[![GitHub Release][release-badge]][release-link]
[![Go Report Card][report-badge]][report-link]
[![License][license-badge]][license-link]

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [helm-repo-updater](#helm-repo-updater)
  - [Scope](#scope)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Examples of usage](#examples-of-usage)
    - [Using the binary](#using-the-binary)
    - [Using a Docker Container](#using-a-docker-container)
  - [Running the tests](#running-the-tests)
    - [Tests requirements](#tests-requirements)
    - [Launch tests](#launch-tests)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# helm-repo-updater

## Scope

This repo aims to manage the development of `helm-repo-updater`, a CLI tool whose objective is to update one or more sets of value keys in the `values.yaml` files used by a [helm chart](https://helm.sh/docs/topics/charts/) stored in a specific git repository, using a commit with the desired change in the repository where the `values.yaml` file is stored.

## Installation

- Launch directly the binary, for that option it will be necessary go to [release page](https://github.com/DocPlanner/helm-repo-updater/releases) and download the binary needed for the architecture of the machine where it is going to run

- It is possible to use a Docker container, since for each [release](https://github.com/DocPlanner/helm-repo-updater/releases) the associated image will be published in the [GitHub Container Registry of the repository](ghcr.io/docplanner/helm-repo-updater). The `latest` tag is also available and will point to the latest published version.

## Usage

    Runs the helm repo updater

    Usage:
      helm-repo-updater run [flags]

    Flags:
          --allow-nothing-to-update          allow the error message 'nothing to update, skipping commit' and finish without exit 1 the execution
          --app-name string                  app name
          --dry-run                          run in dry-run mode. If set to true, do not perform any changes
          --git-branch string                git repo branch (default "develop")
          --git-commit-email string          e-mail address to use for Git commits
          --git-commit-user string           Username to use for Git commits
          --git-dir string                   file eg. /production/charts/
          --git-file string                  file eg. values.yaml
          --git-password string              Password for github user
          --git-repo-url string              git repo url
          --helm-key-values stringToString   helm key-values sets (default [])
      -h, --help                             help for run
          --logLevel string                  set the loglevel to one of trace|debug|info|warn|error (default "info")
          --ssh-private-key string           ssh private key
          --use-ssh-private-key-as-inline    ssh private key inline creation, if true it will use ssh-private-key as input for create ssh private key file in temporal directory

    Global Flags:
          --config string   config file (default is $HOME/.helm-repo-updater.yaml)

## Examples of usage

### Using the binary

- Example run to update the `.image.tag` key to `1.1.0` in the `develop` branch of the `test-repo` repository:
  ```bash
  $ helm-repo-updater run \
    --app-name=example-app \
    --git-branch="develop" \
    --git-commit-user="test-user" \
    --git-commit-email="test-user@docplanner.com" \
    --git-file="values.yaml" \
    --helm-key-values=".image.tag=1.1.0" \
    --git-repo-url="git@github.com:DocPlanner/example-repo.git" \
    --ssh-private-key="test-git-server/private_keys/helm-repo-updater-test"
  INFO[2022-03-03T16:23:10+01:00] Cloning git repository git@github.com:DocPlanner/example-repo.git in temporal folder located in /var/folders/vb/v4wr_9f52ns4mmdkwp4_35cm0000gp/T/git-example-app956297090  application=example-app
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
    --git-repo-url="git@github.com:DocPlanner/example-repo.git" \
    --ssh-private-key="test-git-server/private_keys/helm-repo-updater-test"
  INFO[2022-03-03T16:24:17+01:00] Cloning git repository git@github.com:DocPlanner/example-repo.git in temporal folder located in /var/folders/vb/v4wr_9f52ns4mmdkwp4_35cm0000gp/T/git-example-app1208822616  application=example-app
  Enumerating objects: 8, done.
  Counting objects: 100% (8/8), done.
  Compressing objects: 100% (2/2), done.
  Total 8 (delta 0), reused 0 (delta 0), pack-reused 0
  INFO[2022-03-03T16:24:17+01:00] Pulling latest changes of branch develop      application=example-app
  INFO[2022-03-03T16:24:17+01:00] Actual value for key .image.tag: 1.1.0        application=example-app
  INFO[2022-03-03T16:24:17+01:00] Setting new value for key .image.tag: 1.1.0   application=example-app
  INFO[2022-03-03T16:24:17+01:00] target for key .image.tag is the same, skipping  application=example-app
  ERRO[2022-03-03T16:24:17+01:00] Could not update application spec: nothing to update, skipping commit  application=example-app
  INFO[2022-03-03T16:24:17+01:00] nothing to update, skipping commit  application=example-app
  ```

  - Example run to update the `.image.tag` key to `1.1.0` in the `develop` branch of the `test-repo` repository, being `1.1.0` the value currently present in the repository for the above key without allowing the error with message `nothing to update, skipping commit`:
  ```
  $ helm-repo-updater run \
    --allow-nothing-to-update=false \
    --app-name=example-app \
    --git-branch="develop" \
    --git-commit-user="test-user" \
    --git-commit-email="test-user@docplanner.com" \
    --git-file="values.yaml" \
    --helm-key-values=".image.tag=1.1.0" \
    --git-repo-url="git@github.com:DocPlanner/example-repo.git" \
    --ssh-private-key="test-git-server/private_keys/helm-repo-updater-test"
  INFO[2022-03-03T16:24:17+01:00] Cloning git repository git@github.com:DocPlanner/example-repo.git in temporal folder located in /var/folders/vb/v4wr_9f52ns4mmdkwp4_35cm0000gp/T/git-example-app1208822616  application=example-app
  Enumerating objects: 8, done.
  Counting objects: 100% (8/8), done.
  Compressing objects: 100% (2/2), done.
  Total 8 (delta 0), reused 0 (delta 0), pack-reused 0
  INFO[2022-03-03T16:24:17+01:00] Pulling latest changes of branch develop      application=example-app
  INFO[2022-03-03T16:24:17+01:00] Actual value for key .image.tag: 1.1.0        application=example-app
  INFO[2022-03-03T16:24:17+01:00] Setting new value for key .image.tag: 1.1.0   application=example-app
  INFO[2022-03-03T16:24:17+01:00] target for key .image.tag is the same, skipping  application=example-app
  ERRO[2022-03-03T16:24:17+01:00] Could not update application spec: nothing to update, skipping commit  application=example-app
    ```

### Using a Docker Container

- Example run to update the `.image.tag` key to `1.1.0` in the `develop` branch of the `k8s-argocd-apps` repository for the application `example-app`:
  ```bash
  $ docker run -v ~/.ssh/repositories_keys/example-user_github:/tmp/ssh_key ghcr.io/docplanner/helm-repo-updater run \
      --app-name=example-app \
      --git-branch="develop" \
      --git-commit-user="test-user" \
      --git-commit-email="test-user@docplanner.com" \
      --git-dir="apps/" \
      --git-file="helm/t0/testing/image.yaml" \
      --helm-key-values=".image.tag=1.1.0" \
      --git-repo-url="git@github.com:DocPlanner/example-repo.git" \
      --logLevel=debug \
      --ssh-private-key="/tmp/ssh_key"
  time="2022-03-04T11:29:07Z" level=debug msg="Successfully parsed commit message template" application=example-app
  time="2022-03-04T11:29:07Z" level=debug msg="Processing application example-app in directory apps/example-app/helm/t0/testing/image.yaml"
  time="2022-03-04T11:29:07Z" level=debug msg="Created temporal directory /tmp/git-example-app3237525166 to clone repository git@github.com:DocPlanner/example-repo.git" application=example-app
  time="2022-03-04T11:29:07Z" level=info msg="Cloning git repository git@github.com:DocPlanner/example-repo.git in temporal folder located in /tmp/git-example-app3237525166" application=example-app
  Enumerating objects: 192, done.
  Counting objects: 100% (183/183), done.
  Compressing objects: 100% (99/99), done.
  Total 192 (delta 35), reused 153 (delta 25), pack-reused 9
  time="2022-03-04T11:29:09Z" level=info msg="Pulling latest changes of branch develop" application=example-app
  time="2022-03-04T11:29:10Z" level=info msg="Actual value for key .image.tag: 1.0.0" application=example-app
  time="2022-03-04T11:29:10Z" level=info msg="Setting new value for key .image.tag: 1.1.0" application=example-app
  time="2022-03-04T11:29:10Z" level=debug msg="templated commit message successfully with value: 🚀 automatic update of example-app\nupdates key .image.tag value from '1.0.0' to '1.1.0'\n" application=example-app
  time="2022-03-04T11:29:10Z" level=info msg="Adding file apps/example-app/helm/t0/testing/image.yaml to git for commit changes" application=example-app
  time="2022-03-04T11:29:10Z" level=debug msg="Obtaining current status after changes" application=example-app
  time="2022-03-04T11:29:10Z" level=debug msg="Obtained git status status is: M  apps/example-app/helm/t0/testing/image.yaml\n" application=example-app
  time="2022-03-04T11:29:10Z" level=info msg="It's going to commit changes with message: 🚀 automatic update of example-app\nupdates key .image.tag value from '1.0.0' to '1.1.0'\n" application=example-app
  time="2022-03-04T11:29:10Z" level=debug msg="Obtaining current HEAD to verify added changes" application=example-app
  time="2022-03-04T11:29:10Z" level=info msg="It's going to push commit with hash cb89e6e2a9a238cd0d1d5cdf29cf408185545ed6 and message 🚀 automatic update of example-app\nupdates key .image.tag value from '1.0.0' to '1.1.0'\n" application=example-app
  time="2022-03-04T11:29:10Z" level=info msg="Pushing changes" application=example-app
  time="2022-03-04T11:29:12Z" level=info msg="Successfully pushed changes" application=example-app
  time="2022-03-04T11:29:12Z" level=info msg="Successfully updated the live application spec" application=example-app
  ```

- Example run to update the `.image.tag` key to `1.1.0` in the `develop` branch of the `k8s-argocd-apps` repository for the `example-app`, being `1.1.0` the value currently present in the repository for the above key:
  ```
  $ docker run -v ~/.ssh/repositories_keys/example-user_github:/tmp/ssh_key ghcr.io/docplanner/helm-repo-updater run \
      --app-name=example-app \
      --git-branch="develop" \
      --git-commit-user="test-user" \
      --git-commit-email="test-user@docplanner.com" \
      --git-dir="apps/" \
      --git-file="helm/t0/testing/image.yaml" \
      --helm-key-values=".image.tag=1.1.0" \
      --git-repo-url="git@github.com:DocPlanner/example-repo.git" \
      --logLevel=debug \
      --ssh-private-key="/tmp/ssh_key"
  time="2022-03-04T11:30:48Z" level=debug msg="Successfully parsed commit message template" application=example-app
  time="2022-03-04T11:30:48Z" level=debug msg="Processing application example-app in directory apps/example-app/helm/t0/testing/image.yaml"
  time="2022-03-04T11:30:48Z" level=debug msg="Created temporal directory /tmp/git-example-app2386521755 to clone repository git@github.com:DocPlanner/example-repo.git" application=example-app
  time="2022-03-04T11:30:48Z" level=info msg="Cloning git repository git@github.com:DocPlanner/example-repo.git in temporal folder located in /tmp/git-example-app2386521755" application=example-app
  Enumerating objects: 200, done.
  Counting objects: 100% (191/191), done.
  Compressing objects: 100% (103/103), done.
  Total 200 (delta 37), reused 159 (delta 25), pack-reused 9
  time="2022-03-04T11:30:49Z" level=info msg="Pulling latest changes of branch develop" application=example-app
  time="2022-03-04T11:30:50Z" level=info msg="Actual value for key .image.tag: 1.1.0" application=example-app
  time="2022-03-04T11:30:50Z" level=info msg="Setting new value for key .image.tag: 1.1.0" application=example-app
  time="2022-03-04T11:30:50Z" level=info msg="target for key .image.tag is the same, skipping" application=example-app
  time="2022-03-04T11:30:50Z" level=error msg="Could not update application spec: nothing to update, skipping commit" application=example-app
  time="2022-03-04T11:30:50Z" level=info msg="nothing to update, skipping commit"  application=example-app
  ```

  - Example run to update the `.image.tag` key to `1.1.0` in the `develop` branch of the `k8s-argocd-apps` repository for the `example-app`, being `1.1.0` the value currently present in the repository for the above key without allowing the error with message `nothing to update, skipping commit`::
  ```
  $ docker run -v ~/.ssh/repositories_keys/example-user_github:/tmp/ssh_key ghcr.io/docplanner/helm-repo-updater run \
      --app-name=example-app \
      --git-branch="develop" \
      --git-commit-user="test-user" \
      --git-commit-email="test-user@docplanner.com" \
      --git-dir="apps/" \
      --git-file="helm/t0/testing/image.yaml" \
      --helm-key-values=".image.tag=1.1.0" \
      --git-repo-url="git@github.com:DocPlanner/example-repo.git" \
      --logLevel=debug \
      --ssh-private-key="/tmp/ssh_key"
  time="2022-03-04T11:30:48Z" level=debug msg="Successfully parsed commit message template" application=example-app
  time="2022-03-04T11:30:48Z" level=debug msg="Processing application example-app in directory apps/example-app/helm/t0/testing/image.yaml"
  time="2022-03-04T11:30:48Z" level=debug msg="Created temporal directory /tmp/git-example-app2386521755 to clone repository git@github.com:DocPlanner/example-repo.git" application=example-app
  time="2022-03-04T11:30:48Z" level=info msg="Cloning git repository git@github.com:DocPlanner/example-repo.git in temporal folder located in /tmp/git-example-app2386521755" application=example-app
  Enumerating objects: 200, done.
  Counting objects: 100% (191/191), done.
  Compressing objects: 100% (103/103), done.
  Total 200 (delta 37), reused 159 (delta 25), pack-reused 9
  time="2022-03-04T11:30:49Z" level=info msg="Pulling latest changes of branch develop" application=example-app
  time="2022-03-04T11:30:50Z" level=info msg="Actual value for key .image.tag: 1.1.0" application=example-app
  time="2022-03-04T11:30:50Z" level=info msg="Setting new value for key .image.tag: 1.1.0" application=example-app
  time="2022-03-04T11:30:50Z" level=info msg="target for key .image.tag is the same, skipping" application=example-app
  time="2022-03-04T11:30:50Z" level=error msg="Could not update application spec: nothing to update, skipping commit" application=example-app
  time="2022-03-04T11:30:50Z" level=error msg="Error trying to update the example-app application: nothing to update, skipping commit" application=example-app
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

<!-- JUST BADGES & LINKS -->
[tests-badge]: https://img.shields.io/github/workflow/status/DocPlanner/helm-repo-updater/Tests
[tests-link]: https://github.com/DocPlanner/helm-repo-updater/actions?query=workflow%3ATests

[release-badge]: https://img.shields.io/github/release/DocPlanner/helm-repo-updater.svg?logo=github&labelColor=262b30
[release-link]: https://github.com/DocPlanner/helm-repo-updater/releases

[report-badge]: https://goreportcard.com/badge/github.com/DocPlanner/helm-repo-updater
[report-link]: https://goreportcard.com/report/github.com/DocPlanner/helm-repo-updater

[license-badge]: https://img.shields.io/github/license/DocPlanner/helm-repo-updater
[license-link]: https://github.com/DocPlanner/helm-repo-updater/blob/master/LICENSE
