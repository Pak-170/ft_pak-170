#!/bin/sh
set -xe

cd /app
go mod download
air -c .air.toml