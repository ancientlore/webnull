#!/bin/bash

GO_VERSION=$(go env GOVERSION | cut -b "3-")
GO_MAJOR_VERSION=$(cut -d '.' -f 1,2 <<< "$GO_VERSION")
TAG=$(git tag | sort -V | tail -1)

echo
echo Go version is $GO_VERSION, major version is $GO_MAJOR_VERSION
echo Tag is $TAG

echo
echo Building ancientlore/webnull:$TAG
docker build --build-arg GO_VERSION=$GO_VERSION --build-arg IMG_VERSION=$GO_MAJOR_VERSION -t ancientlore/webnull:$TAG . || return 1

echo
echo Pushing ancientlore/webnull:$TAG
docker push ancientlore/webnull:$TAG

echo
echo Tagging ancientlore/webnull:latest
docker tag ancientlore/webnull:$TAG ancientlore/webnull:latest

echo
echo Pushing ancientlore/webull:latest
docker push ancientlore/webnull:latest

echo
echo Tagging ancientlore.registry.cpln.io/webnull:$TAG
docker tag ancientlore/webnull:$TAG ancientlore.registry.cpln.io/webnull:$TAG

echo
echo Pushing ancientlore.registry.cpln.io/webull:$TAG
docker push ancientlore.registry.cpln.io/webnull:$TAG
