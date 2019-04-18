#!/bin/sh

cd cmd/vivograph
go build
cd ../../
cd cmd/vivograph-cli
go build


