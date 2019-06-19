#!/bin/bash

set -ex

cd ./selfhydro
go get .
go test -cover ./... | tee test_coverage.txt

mv test_coverage.txt $GOPATH/coverage-results/.
