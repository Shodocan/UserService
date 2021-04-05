#!/bin/bash

for file in $(ls deployments/k8s/); do kubectl apply -f deployments/k8s/$file; done