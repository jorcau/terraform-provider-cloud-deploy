#!/bin/sh

PACKAGE_PATH="${GOPATH}/src/cloud-deploy.io/${BITBUCKET_REPO_SLUG}"
mkdir -pv "${PACKAGE_PATH}"
tar -cO --exclude-vcs --exclude=bitbucket-pipelines.yml . | tar -xv -C "${PACKAGE_PATH}"
cd "${PACKAGE_PATH}"
git ls-remote --heads --tags https://github.com/jorcau/terraform-provider-cloud-deploy.git | grep -E "refs/(heads|tags)/${BITBUCKET_TAG}$"
curl -sL https://git.io/goreleaser | bash