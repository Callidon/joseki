#!/bin/bash
go test -coverprofile=graph.cover.out -coverpkg=./... graph
go test -coverprofile=graph.cover.out -coverpkg=./... ./graph
go test -coverprofile=parser.cover.out -coverpkg=./... ./parser
go test -coverprofile=rdf.cover.out -coverpkg=./... ./rdf
go test -coverprofile=sparql.cover.out -coverpkg=./... ./sparql
echo "mode: set" > coverage.out && cat *.cover.out | grep -v mode: | sort -r | awk '{if($1 != last) {print $0;last=$1}}' >> coverage.out
