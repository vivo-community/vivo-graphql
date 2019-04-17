#!/bin/sh

cd cmd/vivoql
go build
cd ../../
cd cmd/vivo_mappings
go build
cd ../../

