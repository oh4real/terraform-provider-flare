#!/bin/bash -e

OSs=("darwin" "linux")
ARCHs=("386" "amd64")

# export GOPATH="$HOME/go"

# create build dir location
rootDir=`pwd`
releases="${rootDir}/releases"
mkdir $releases

#Get into the right directory
cd $(dirname $0)

#Parse command line params
CONFIG=$@
for line in $CONFIG; do
  eval "$line"
done

if [[ -z "$github_api_token" && -f github_api_token ]];then
  github_api_token=$(cat github_api_token)
fi

if [[ -z "$owner" ]];then
  owner="oh4real"
fi

if [[ -z "$repo" ]];then
  repo="terraform-provider-flare"
fi

if [[ -z "$github_api_token" || -z "$owner" || -z "$repo" || -z "$tag" ]];then
  echo "USAGE: $0 github_api_token=TOKEN owner=someone repo=somerepo tag=vX.Y.Z"
  exit 1
fi

if [[ "$tag" != v* ]];then
  tag="v$tag"
fi

#Build for all architectures we want
ARTIFACTS=()
#for GOOS in darwin linux windows netbsd openbsd solaris;do

echo "Building..."
for GOOS in "${OSs[@]}";do
  for GOARCH in "${ARCHs[@]}";do
    export GOOS GOARCH

    TF_OUT_FILE="$releases/terraform-provider-flare-$GOOS-$GOARCH"
    echo "  $TF_OUT_FILE"
    go build -o "$TF_OUT_FILE" ../
    ARTIFACTS+=("$TF_OUT_FILE")

  done
done

#Create the release so we can add our files
./create-github-release.sh github_api_token=$github_api_token owner=$owner repo=$repo tag=$tag draft=false

#Upload all of the files to the release
for FILE in "${ARTIFACTS[@]}";do
  ./upload-github-release-asset.sh github_api_token=$github_api_token owner=$owner repo=$repo tag=$tag filename="$FILE"
done

echo "Cleaning up..."
rm -f release_info.md
rm -rf $releases
