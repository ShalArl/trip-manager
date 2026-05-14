# Deployment Guide

This guide provides instructions for deploying the application using Helm charts. Helm is a package manager for Kubernetes that simplifies the deployment and management of applications.

## Prerequisites
- Kubernetes cluster (e.g., Minikube, GKE, EKS)
- Helm installed and configured to access your Kubernetes cluster

## Secret Management
We use external-secrets to manage sensitive information

### Setup
1. Install External Secrets Operator (Note: This is a separate namespace)
   ```bash
   helm repo add external-secrets https://charts.external-secrets.io
   helm repo update
   helm install external-secrets external-secrets/external-secrets \
     --namespace external-secrets --create-namespace
   ```
   
2. Create a secret in gcloud using the following command:
   ```bash
   gcloud secrets create <secret-name> --data-file=<path-to-secret-file>
   ```
   
   or alternatively for one-liners:
   ```bash
    echo -n "<secret-value>" | gcloud secrets create <secret-name> --data-file=-
   ```

## Deployment Steps (local development)
1. Update Chart:
   ```bash
   helm dependency update helm/services/<service-name>
   ```
2. Prepare image
   Build the Docker image and tag it and make sure it is loaded into your Kubernetes cluster. For local development, you can build the image locally and load it into Minikube or your local Kubernetes cluster.
   ```bash
   eval $(minikube docker-env)
   docker build -f backend/<service-name>/Dockerfile backend/ -t localhost/<service-name>:latest
   ```
3. Create a kubernetes namespace (optional - skip if already exists):
   ```bash
   kubectl create namespace <namespace-name>
   ```
4. Install the Helm chart:
   ```bash
   helm install <release-name> helm/service/<service-name> \
     --set image.repository=localhost/<service-name> \
     --set image.tag=latest \
     --values helm/services/<service-name>/values.yaml \
     --values helm/services/<service-name>/values/dev.yaml \
     --namespace <namespace-name>
   ```
5. Validate the deployment:
   ```bash
   kubectl get pods -n <namespace-name>
   ```
   

## Updating Charts
To update the Helm chart for a service, follow these steps:
1. Make changes to the Helm chart templates or values files as needed.
2. Update the Docker image if necessary and ensure it is available in your Kubernetes cluster.
3. Run the following command to apply the changes:
   ```bash
   helm upgrade <release-name> helm/service/<service-name> \
     --set image.repository=localhost/<service-name> \
     --set image.tag=latest \
     --values helm/services/<service-name>/values.yaml \
     --values helm/services/<service-name>/values/dev.yaml \
     --namespace <namespace-name>
   ```
   For services without a dedicated docker image such as `gateway` for example, use:
   ```bash
   helm upgrade <service-name> helm/services/<service-name> --namespace <namespace-name> --values helm/services/<service-name>/values.yaml --values helm/services/<service-name>/values/dev.yaml
   ```
   

## Something unfortunate happened?

Use the following command to check the logs of the deployed service:
```bash
kubectl logs -n <namespace-name> <pod-name>
```

If the pod is in a restart loop, check the logs of the previous instance:
```bash
kubectl logs -n <namespace-name> <pod-name> --previous
```

Use the following command to get more details about the pod status:
```bash
kubectl describe pod -n <namespace-name> <pod-name>
```

Deployment successful but the service doesn't show up when executing 'kubectl get pods -n <namespace-name>'?
- Validate the deployment by executing: `helm get manifest <service-name> -n <namespace-name>`. This will show you the generated Kubernetes manifests. Check if the resources are created as expected.


## Changing something in the code?

If the helm chart is unchanged and you want to deploy code changes, a helm upgrade won't help because no changes are detected (locally) due to tag latest being used.
In this case you have to manually restart the pod as follows:
```bash
kubectl rollout restart deployment/<deployment-name> -n <namespace-name>
```