#!/bin/bash


#THIS SHELL SCRIPT IS ESPECIALLY FOR DOCKER CONTAINER.
./datahub --daemon --token $DAEMON_TOKEN

while [ 1 ]
do
    sleep 10
done
