#!/usr/bin/env bash

APP_ROOT=$(dirname $0)/..
cd ${APP_ROOT}

VER=${VER:-local-$(date +%Y%m%d%H%M)}
PROJECT=akiho-playground

gcloud config set project ${PROJECT}

export API_IMAGE=gcr.io/${PROJECT}/api:${VER}
docker build . -t ${API_IMAGE} --build-arg _ENTRYPOINT=entrypoint/api/main.go --target deploy
docker login -u oauth2accesstoken -p "$(gcloud auth print-access-token)" https://gcr.io
docker push ${API_IMAGE}

gcloud container clusters get-credentials app-cluster --zone=asia-northeast1-a
envsubst < k8s/api.yaml | cat | kubectl apply -f -

docker rmi -f `docker images | grep "gcr.io/akiho-playground" | awk '{print $3}'`
docker rmi -f `docker images | grep "<none>" | awk '{print $3}'`