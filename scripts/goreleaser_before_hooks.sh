#!/bin/sh

git ls-remote --heads --tags https://github.com/jorcau/terraform-provider-cloud-deploy.git | grep -E "refs/(heads|tags)/${version}$"
