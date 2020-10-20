#!/bin/sh -ex

if [ -f artifact_go.sh ]; then
  rm -f artifact_go.sh
fi
wget -q -O artifact_go.sh 'https://raw.githubusercontent.com/mdblp/tools/feature/add-platform-project/artifact/artifact_go.sh'
chmod +x artifact_go.sh

export ARTIFACT_GO_VERSION='1.11.4'

chmod 755 build.sh
# Disable binary deployment in artifact_go.sh as it is done by the Make command
export ARTIFACT_DEPLOY=false
export ARTIFACT_BUILD=true
export BUILD_OPENAPI_DOC=false
export SECURITY_SCAN=true
./artifact_go.sh data 
