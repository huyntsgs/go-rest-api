#!/bin/bash

cd go-rest-api

go build

cd migration

./migration

cd ..

./go-rest-api
