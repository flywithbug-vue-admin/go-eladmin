#!/usr/bin/env bash

CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build main.go
scp main root@118.89.108.25:/root/go-eladmin/vue-admin1

rm -r main