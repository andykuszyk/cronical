#!/bin/bash
tag=$(echo -n $GITHUB_REF | sed 's/refs\/tags\///g')
docker build -t andykuszyk/cronical:$tag .
docker push andykuszyk/cronical:$tag
