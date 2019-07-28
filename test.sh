#!/usr/bin/env bash

aws dynamodb create-table \
  --provisioned-throughput '{"ReadCapacityUnits": 5, "WriteCapacityUnits": 5}' \
  --table-name login \
  --attribute-definition '[{ "AttributeName": "identifier", "AttributeType": "S"}]' \
  --key-schema '[{ "KeyType":"HASH", "AttributeName":"identifier"}]' \
  --endpoint-url http://docker.devel:4569

#aws dynamodb put-item \
#  --table-name login \
#  --item '{"identifier": {"S": "5f46cf19-5399-55e3-aa62-0e7c19382250"}, "email": {"S": "tester@carpark.ninja"}}' \
#  --endpoint-url http://docker.devel:4569

#aws dynamodb put-item \
#--table-name login \
#--item '{"identifier": {"S": "test1"}, "validto": {"N": "1539261683"}}' \
#--endpoint-url http://docker.devel:8000
