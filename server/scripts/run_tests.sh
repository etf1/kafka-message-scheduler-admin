#!/usr/bin/env bash

set -euo pipefail

function wait_for_kafka(){
    local server="$1"

    for i in {1..10}
    do
        echo "Waiting for kafka cluster to be ready ..."
        # kafkacat has 5s timeout
        kafkacat -b "${server}" -L > /dev/null 2>&1 && break
    done
}

function wait_for_scheduler(){
    local server="$1"

    for i in {1..10}
    do
        echo "Waiting for scheduler api to be ready ..."
        http_code=$(curl -s -o /dev/null -w "%{http_code}" ${server} || true)
        if [[ "${http_code}" == "200" ]]; then
            break
        fi

        sleep 5
    done
}

wait_for_kafka "kafka:29092"

wait_for_scheduler "scheduler:8000/info"
wait_for_scheduler "scheduler:8000/schedules"

RUN_INTEGRATION_TESTS=yes make tests