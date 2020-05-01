#!/bin/bash
# prints the ip address of the myhttp docker container
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' myhttp
