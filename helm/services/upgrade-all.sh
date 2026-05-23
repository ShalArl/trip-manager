#!/bin/bash

for svc in auth presigner social users trips frontend gateway; do
  helm dependency update helm/services/$svc

  helm upgrade $svc helm/services/$svc \
    --namespace trip-manager-dev \
    --values helm/services/$svc/values.yaml \
    --values helm/services/$svc/values/dev.yaml \
    --set image.repository=localhost/$svc \
    --set image.tag=latest
done