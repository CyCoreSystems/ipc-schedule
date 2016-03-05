#!/bin/bash

go build
./ipc-schedule -debug &
./node_modules/.bin/webpack-dev-server -d --inline --hot --port 3000 
