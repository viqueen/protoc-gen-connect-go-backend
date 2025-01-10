#!/usr/bin/env bash

function local() {
  goreleaser build --clean --snapshot
}

eval "$@"