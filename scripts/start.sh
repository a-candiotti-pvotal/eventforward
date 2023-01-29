#!/bin/sh

env $(cat .env | xargs) go run cmd/eventforward/main.go
