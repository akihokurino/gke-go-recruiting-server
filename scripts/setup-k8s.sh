#!/usr/bin/env bash

APP_ROOT=$(dirname $0)/..

cd ${APP_ROOT}

PROJECT=gke-go-sample-dev
REGION=asia-northeast1
VPC_NAME=gke-go-sample
SUBNET_NAME=gke-go-sample-subnet

gcloud config set project ${PROJECT}

echo -e "\n\033[1;32m----- VPC作成 -----\033[0;39m"
gcloud compute networks create ${VPC_NAME} \
  --project=${PROJECT} \
  --bgp-routing-mode=regional \
  --subnet-mode=custom

echo -e "\n\033[1;32m----- Subnet作成 -----\033[0;39m"
gcloud compute networks subnets create ${SUBNET_NAME} \
  --project=${PROJECT} \
  --region=${REGION} \
  --network=${VPC_NAME} \
  --range=192.168.1.0/24

echo -e "\n\033[1;32m----- IP作成 -----\033[0;39m"
gcloud compute addresses create gke-go-sample --region=${REGION}
gcloud compute addresses create api-ip --global
gcloud beta compute ssl-certificates create api-cert --domains dev-api.gke-go-sample.jp

echo -e "\n\033[1;32m----- Router作成 -----\033[0;39m"
gcloud compute routers create gke-go-sample \
  --region=${REGION} \
  --network=${VPC_NAME} \
  --asn=65001

echo -e "\n\033[1;32m----- Nat作成 -----\033[0;39m"
gcloud compute routers nats create gke-go-sample \
  --region=${REGION} \
  --router=gke-go-sample \
  --nat-external-ip-pool="gke-go-sample" \
  --nat-custom-subnet-ip-ranges="${SUBNET_NAME}"

echo -e "\n\033[1;32m----- GKE Cluster作成 -----\033[0;39m"
gcloud container clusters create api-cluster \
  --project=${PROJECT} \
  --zone=asia-northeast1-a \
  --network=${VPC_NAME} \
  --subnetwork=${SUBNET_NAME} \
  --enable-ip-alias \
  --enable-private-nodes \
  --master-ipv4-cidr=172.16.0.0/28 \
  --enable-master-authorized-networks \
  --master-authorized-networks=0.0.0.0/0 \
  --no-enable-legacy-authorization \
  --no-enable-basic-auth \
  --no-issue-client-certificate \
  --machine-type=e2-micro \
  --disk-size=30 \
  --num-nodes=1

gcloud container clusters create batch-cluster \
  --project=${PROJECT} \
  --zone=asia-northeast1-c \
  --network=${VPC_NAME} \
  --subnetwork=${SUBNET_NAME} \
  --enable-ip-alias \
  --enable-private-nodes \
  --master-ipv4-cidr=172.16.1.0/28 \
  --enable-master-authorized-networks \
  --master-authorized-networks=0.0.0.0/0 \
  --no-enable-legacy-authorization \
  --no-enable-basic-auth \
  --no-issue-client-certificate \
  --machine-type=e2-micro \
  --disk-size=30 \
  --num-nodes=1