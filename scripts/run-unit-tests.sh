#!/bin/bash

printf "Running Unit Tests\n"

function print_success {
    printf "\n"
    printf '\e[1;32m%-6s\e[m\n' "******************************************************"
    printf '\e[1;32m%-6s\e[m\n' "All tests passed."
    printf '\e[1;32m%-6s\e[m\n' "******************************************************"
}

function print_failure {
    printf "\n"
    printf '\e[1;31m%-6s\e[m\n' "******************************************************"
    printf '\e[1;31m%-6s\e[m\n' "One or more tests failed."
    printf '\e[1;31m%-6s\e[m\n' "******************************************************"
}

if [ "$1" == "-coverage" ]; then
    echo "Running it with Code Coverage"
    gotestsum --format testname --junitfile unit-tests.xml -- -race -covermode=atomic -coverprofile=coverage.out ./...
    test_status=$?  
    if [ "$2" == "-upload" ]; then
        echo "Moving coverage report and junit-xml in tmp folder"
        go tool cover -html=coverage.out -o coverage.html;
        mv coverage.html ./tmp/artifacts
    else 
        echo "Opening up coverage report in browser"
        sleep 2 && go tool cover -html=coverage.out;
    fi
else
    gotestsum --format testname ./...
    test_status=$?  
fi

if [ $test_status == 0 ]; 
then
    print_success
else
    print_failure
fi
exit $test_status
