#!/bin/bash
set -o xtrace
docker build -t ming/chilli .
docker run -d -p 80:8080 ming/chilli
