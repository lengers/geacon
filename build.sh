#!/bin/bash

# copy the .cobaltstrike.beacon_keys file to the current directory
docker cp cobaltstrike:/opt/cobaltstrike/.cobaltstrike.beacon_keys .

# extract the keys
export PRIVATE_KEY=$(java --enable-preview -jar tools/BeaconTool/BeaconTool.jar -i .cobaltstrike.beacon_keys -rsa | perl -0777 -ne '/(-----BEGIN PRIVATE KEY-----.+?-----END PRIVATE KEY-----)/sg && print"$1\n"')
export PUBLIC_KEY=$(java --enable-preview -jar tools/BeaconTool/BeaconTool.jar -i .cobaltstrike.beacon_keys -rsa | perl -0777 -ne '/(-----BEGIN PUBLIC KEY-----.+?-----END PUBLIC KEY-----)/sg && print"$1\n"')

# replace the keys in the code
perl -0777 -i -pe 's/-----BEGIN PRIVATE KEY-----.+?-----END PRIVATE KEY-----/$ENV{"PRIVATE_KEY"}/gs' cmd/config/config.go
perl -0777 -i -pe 's/-----BEGIN PUBLIC KEY-----.+?-----END PUBLIC KEY-----/$ENV{"PUBLIC_KEY"}/gs' cmd/config/config.go

# build the program
docker run -ti --rm -v $PWD:/usr/src/app golang:buster bash -c \
    'export GOOS="linux" && \
    export GOARCH="amd64" && \
    cd /usr/src/app && \
    go build cmd/main.go '
