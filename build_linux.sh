#!/bin/bash
set -e

(cd cmd && (
	ls *.go | while read i ; do
		GOOS=linux GOARCH=amd64 go build -o `dirname $i`/`basename $i .go`_linux_amd64 $i
	done
))

(cd scripts && (
	ls *.go | while read i ; do
		GOOS=linux GOARCH=amd64 go build -o `dirname $i`/`basename $i .go`_linux_amd64 $i
	done
))
