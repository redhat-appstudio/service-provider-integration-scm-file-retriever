#!/usr/bin/env bash
echo 'restarting minikube'
minikube stop
minikube delete
set -e
minikube start --cpus 4 --memory 6000
minikube addons enable ingress
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.6.1/cert-manager.yaml


kubectl rollout status deployment/cert-manager  -n cert-manager
kubectl rollout status deployment/cert-manager-cainjector  -n cert-manager
kubectl rollout status deployment/cert-manager-webhook  -n cert-manager

kubectl rollout status deployment/ingress-nginx-controller  -n ingress-nginx
