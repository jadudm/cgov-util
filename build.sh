#!/bin/bash

# go build -ldflags "-X main.SHA1=${{steps.vars.outputs.sha_short}}" -v ./...
go build -ldflags "-X main.SHA1=$sha_short" .