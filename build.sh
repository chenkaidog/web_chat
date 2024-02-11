#!/bin/bash
RUN_NAME=hertz_service
mkdir -p output/bin output/conf output/static
cp script/* output/
cp -r conf/ output/
cp -r static/ output/
chmod +x output/bootstrap.sh

go build -o output/bin/${RUN_NAME}