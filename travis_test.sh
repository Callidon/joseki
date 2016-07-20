#!/bin/bash
PACKAGES="graph parser rdf sparql"
for pkg in $PACKAGES; do
  go test -coverprofile=$pkg.cover.out -coverpkg=./... ./$pkg
done
echo "mode: set" > coverage.out && cat *.cover.out | grep -v mode: | sort -r | awk '{if($1 != last) {print $0;last=$1}}' >> coverage.out
