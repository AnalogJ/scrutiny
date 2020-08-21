#! /bin/sh

export GO_PACKAGES=$(go list ./... | grep -v /vendor/)
