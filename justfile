#!/usr/bin/env just --justfile

set shell := ["bash", "-c"]

test PKG="./..." *ARGS="":
    gotestsum -- -race -count 1 -coverprofile cover.out {{ ARGS }} {{ PKG }}

vtest PKG="./..." *ARGS="":
    gotestsum -f testname -- -race -count 1 -coverprofile cover.out {{ ARGS }} {{ PKG }}

bench PKG="./..." *ARGS="":
    go test -v -count 1 {{ PKG }} -run "no tests" -bench . {{ ARGS }}

lint PKG="./..." *ARGS="--new=false --fix":
    golangci-lint run {{ PKG }} {{ ARGS }}
