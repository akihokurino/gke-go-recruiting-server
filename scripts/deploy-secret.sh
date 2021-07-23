#!/usr/bin/env bash

APP_ROOT=$(dirname $0)/..
cd ${APP_ROOT}

VER=${VER:-local-$(date +%Y%m%d%H%M)}
PROJECT=gke-go-sample-dev

gcloud config set project ${PROJECT}

echo "--------------- start update secret ---------------"

gcloud container clusters get-credentials api-cluster --zone=asia-northeast1-a
kubectl delete secret env
kubectl create secret generic gcp-credentials --from-file=credentials.json=./config/gcp-key.json
kubectl create secret generic firebase-credentials --from-file=credentials.json=./config/firebase-key.json
kubectl create secret generic env --from-file=env=${APP_ROOT}/config/.env

gcloud container clusters get-credentials batch-cluster --zone=asia-northeast1-c
kubectl delete secret env
kubectl create secret generic gcp-credentials --from-file=credentials.json=./config/gcp-key.json
kubectl create secret generic firebase-credentials --from-file=credentials.json=./config/firebase-key.json
kubectl create secret generic env --from-file=env=${APP_ROOT}/config/.env

echo "--------------- complete update secret ---------------"