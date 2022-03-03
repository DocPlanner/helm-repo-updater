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
    - [Example of usage](#example-of-usage)
  - [pre-commit](#pre-commit)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# helm-repo-updater

## Scope

This repo aims to manage the development of `helm-repo-updater`, a CLI tool whose objective is to update one or more sets of value keys in the `values.yaml` files used by a [helm chart](https://helm.sh/docs/topics/charts/) stored in a specific git repository, using a commit with the desired change in the repository where the `values.yaml` file is stored.

## Installation

Go to [release page](https://github.com/DocPlanner/helm-repo-updater/releases) and download the binary needed for the architecture of the machine where it is going to run

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

### Example of usage

- Example run to update the `.image.tag` key to `1.1.0` in the `develop` branch of the `test-repo` repository::
  ```bash
  $ helm-repo-updater run --app-name=example-app --git-branch="develop" --git-commit-user="test-user" --git-commit-email="test-user@docplanner.com" --git-file="values.yaml" --helm-key-values=".image.tag=1.1.0" --git-repo-url="ssh://git@localhost:2222/git-server/repos/test-repo.git" --ssh-private-key="test-git-server/private_keys/helm-repo-updater-test"
  INFO[2022-03-03T15:25:53+01:00] Cloning git repository ssh://git@localhost:2222/git-server/repos/test-repo.git in temporal folder located in /var/folders/vb/v4wr_9f52ns4mmdkwp4_35cm0000gp/T/git-example-app4219399441  application=example-app
  Enumerating objects: 4, done.
  Counting objects: 100% (4/4), done.
  Total 4 (delta 0), reused 0 (delta 0), pack-reused 0
  INFO[2022-03-03T15:25:53+01:00] Pulling latest changes of branch develop      application=example-app
  INFO[2022-03-03T15:25:53+01:00] Actual value for key .image.tag: 1.0.0        application=example-app
  INFO[2022-03-03T15:25:53+01:00] Setting new value for key .image.tag: 1.1.0   application=example-app
  INFO[2022-03-03T15:25:53+01:00] Adding file example-app/values.yaml to git for commit changes  application=example-app
  INFO[2022-03-03T15:25:53+01:00] It's going to commit changes with message: /var/folders/vb/v4wr_9f52ns4mmdkwp4_35cm0000gp/T/example-app1005856236  application=example-app
  INFO[2022-03-03T15:25:53+01:00] It's going to push commit with hash 7b690b0d286dac13c23418807a0612acf7a83867 and message /var/folders/vb/v4wr_9f52ns4mmdkwp4_35cm0000gp/T/example-app1005856236  application=example-app
  INFO[2022-03-03T15:25:53+01:00] Pushing changes                               application=example-app
  INFO[2022-03-03T15:25:53+01:00] Successfully pushed changes                   application=example-app
  INFO[2022-03-03T15:25:53+01:00] Successfully updated the live application spec  application=example-app
  ```

## pre-commit

This repo leverage [pre-commit](https://pre-commit.com) to lint, secure, document the codebase. The [pre-commit](https://pre-commit.com) configuration require the following dependencies installed locally:
- [pre-commit](https://pre-commit.com/#install)
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

**One first repo download, to install the pre-commit hooks it's necessary execute the following command**:
```
pre-commit install
```

**To run the hooks at will it's necessary execute the following command**:
```
pre-commit run -a
```
