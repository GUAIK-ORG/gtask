#!/bin/bash

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags -w -o ./release/gtask ./main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags -w -o ./release/client ./client/client.go