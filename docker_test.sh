#!/usr/bin/env bash

# create the image if needed
# docker build -t adsholoko/smallgo .

# push the image to docker.io
# docker push adsholoko/smallgo:latest
# run docker locally
set -u
docker run \
    -it \
    --rm \
    -w /tmp/smallgo \
    -v $(pwd):/tmp/smallgo \
    adsholoko/smallgo:latest make test
