#!/usr/bin/env bash
if ! type "localstack" > /dev/null; then
    docker-compose stop db
    docker-clean
    docker ps -a --format '{{.Names}} {{.Status}}' | grep 'Exited' | awk '{print $1}' | xargs docker rm
    docker-compose up -d db
fi
#STACK=$(SERVICES=dynamodb TMPDIR=private$TMPDIR localstack start --docker)

aws dynamodb delete-table \
  --table-name login \
  --endpoint http://docker.devel:4569


aws dynamodb create-table \
  --table-name login \
  --attribute-definitions AttributeName=identifier,AttributeType=S \
  --key-schema AttributeName=identifier,KeyType=HASH \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --endpoint-url http://docker.devel:4569
