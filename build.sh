#!/bin/bash -x

name=ifaddr
for os in darwin linux windows; do
    for arch in amd64 386 arm; do
        GOOS=$os GOARCH=$arch go build
        if [ $os = "windows" ]; then
            zip "${name}_${os}_${arch}.zip" "${name}.exe"
        else
            tar zcf "${name}_${os}_${arch}.tar.gz" "${name}"
            rm -f "${name}"
        fi
    done
done
