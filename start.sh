#!/bin/bash
go build
cd migration
go build
./migration
cd ..
./go-rest-api
