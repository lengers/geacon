docker run -ti --rm -v $PWD:/usr/src/app golang:buster bash -c \
    'export GOOS="linux" && \
    export GOARCH="amd64" && \
    cd /usr/src/app && \
    go build cmd/main.go '
