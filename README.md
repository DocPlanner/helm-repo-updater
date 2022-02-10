<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [helm-repo-updater](#helm-repo-updater)
  - [Scope](#scope)
  - [pre-commit](#pre-commit)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# helm-repo-updater

## Scope

This repo aims to manage the development of `helm-repo-updater`, a CLI tool whose objective is to update one or more sets of value keys in the `values.yaml` files used by a [helm chart](https://helm.sh/docs/topics/charts/) stored in a specific git repository, using a commit with the desired change in the repository where the `values.yaml` file is stored.

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
