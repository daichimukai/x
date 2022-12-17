#!/bin/bash

set -eux

docker_compose_yaml=./tests/docker-compose.yaml

docker-compose -f ${docker_compose_yaml} build --no-cache
docker-compose -f ${docker_compose_yaml} up -d
HOST_2_LOOPBACK_IP=10.100.220.3
docker-compose -f ${docker_compose_yaml} exec -T host1 ping -c 5 ${HOST_2_LOOPBACK_IP}

TEST_RESULT=$?
if [ ${TEST_RESULT} -eq 0 ]; then
	printf "\e[32m%s\e[m\n" "integration test successed"
else
	printf "\e[31m%s\e[m\n" "integration test failed"
fi

exit ${TEST_RESULT}
