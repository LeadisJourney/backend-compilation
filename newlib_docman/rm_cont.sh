#!/bin/bash
# Delete all containers
docker rm $(docker ps -a -q)
# Delete all images
docker rmi $(docker images -q)
# Delete all volumes
rm -rf /var/lib/docker/volumes/*
