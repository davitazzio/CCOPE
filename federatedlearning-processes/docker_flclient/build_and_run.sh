#! /bin/bash

cd ~/federatedlearning-processes/

docker build -t lucaserf/flclient:latest ./docker_flclient/
docker push lucaserf/flclient:latest

kubectl rollout restart deployment/flclient-deployment

# kubectl apply -f ./docker_flclient/app/deploy/flclient_deploy.yaml
# kubectl delete -f ./docker_flclient/app/deploy/flclient_deploy.yaml