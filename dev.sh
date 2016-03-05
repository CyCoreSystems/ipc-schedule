#!/bin/bash

go generate
go build
./ipc-schedule -debug -addr :9000 &
./node_modules/.bin/webpack-dev-server -d --inline --hot --port 3000 
