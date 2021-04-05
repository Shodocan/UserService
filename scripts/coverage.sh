#!/bin/sh

THRESHOLD=$1


sed -i '/mock.go/d' .coverage/coverage.out
COVERAGE=$(go tool cover -func=.coverage/coverage.out | grep total | awk '{print $3}')
COVERAGE=${COVERAGE%\%}
COVERAGE=$(printf "%.0f" ${COVERAGE} ) 

if [ ${COVERAGE} -gt ${THRESHOLD} ]
then
    echo "coverage above threshold"
    echo "coverage: ${COVERAGE} - threshold: ${THRESHOLD}"
    exit 0
fi

echo "coverage below threshold"
echo "coverage: ${COVERAGE} - threshold: ${THRESHOLD}"
exit 1
