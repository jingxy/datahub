#!/bin/bash


#THIS SHELL SCRIPT IS ESPECIALLY FOR DOCKER CONTAINER.
if [ $DAEMON_TOKEN ]; then
    ./datahub --daemon --token $DAEMON_TOKEN
else
    ./datahub --daemon
fi

while [ 1 ]
do
    sleep 10
done
