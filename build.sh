#!/bin/bash

GO_VERSION=$(go env GOVERSION | cut -b "3-")
GO_MAJOR_VERSION=$(cut -d '.' -f 1,2 <<< "$GO_VERSION")
TAG=$(git tag | sort -V | tail -1)

echo
echo Go version is $GO_VERSION, major version is $GO_MAJOR_VERSION
echo Tag is $TAG

echo
echo Building ancientlore/webnull:$TAG
docker buildx build --build-arg GO_VERSION=$GO_VERSION --build-arg IMG_VERSION=$GO_MAJOR_VERSION --platform linux/amd64,linux/arm64 -t ancientlore/webnull:$TAG . || exit 1

gum confirm "Push?" || exit 1

echo
echo Pushing ancientlore/webnull:$TAG
docker push ancientlore/webnull:$TAG || exit 1

echo
echo Tagging ancientlore/webnull:latest
docker tag ancientlore/webnull:$TAG ancientlore/webnull:latest || exit 1

echo
echo Pushing ancientlore/webull:latest
docker push ancientlore/webnull:latest || exit 1

echo
echo Tagging ancientlore.registry.cpln.io/webnull:$TAG
docker tag ancientlore/webnull:$TAG ancientlore.registry.cpln.io/webnull:$TAG || exit 1

echo
echo Pushing ancientlore.registry.cpln.io/webull:$TAG
docker push ancientlore.registry.cpln.io/webnull:$TAG || exit 1
