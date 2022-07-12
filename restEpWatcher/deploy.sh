#!/bin/bash
kubectl apply -f svcAccount.yaml
sleep 2
kubectl apply -f EPdaemonSet.yaml
