#!/bin/bash

ACTION="${1}"

if [[ "$ACTION" == "create" ]]; then
    echo "It's going to launch actions for create tests environment"
    # Create git-server
    docker-compose -f test-git-server/docker-compose.yaml up -d --build git-server
    # Create repo in git-server
    docker-compose -f test-git-server/docker-compose.yaml up --build create-repo
elif [[ "$ACTION" == "destroy" ]]; then
    # Delete git-server and create-repo containers
    docker-compose -f test-git-server/docker-compose.yaml up down
else
    echo "ACTION provided $ACTION is not a valid command. The only options available are 'create' and 'destroy'"
fi