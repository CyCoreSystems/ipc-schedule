#!/bin/bash

go run ./main.go -mode dev &
./node_modules/.bin/webpack-dev-server -d --inline --hot --port 3000 
