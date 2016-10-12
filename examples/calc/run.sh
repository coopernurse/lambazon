#!/bin/sh

set -e

docker run --rm -p 2015:2015 \
       -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
       -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
       -e AWS_REGION=$AWS_REGION \
       -v `pwd`/Caddyfile:/etc/Caddyfile \
       coopernurse/caddy-awslambda
