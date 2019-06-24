#!/bin/bash

set -ex

cd ./BigQueryPublisher
go get .
go test -cover | tee test_coverage.txt

mv test_coverage.txt ../coverage-results/.
