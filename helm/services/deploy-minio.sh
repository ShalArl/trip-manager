#!/bin/bash

cat > /tmp/minio.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: minio
  namespace: trip-manager-dev
spec:
  replicas: 1
  selector:
    matchLabels:
      app: minio
  template:
    metadata:
      labels:
        app: minio
    spec:
      containers:
        - name: minio
          image: quay.io/minio/minio:latest
          command:
            - minio
            - server
            - /data
            - --console-address
            - ":9001"
          ports:
            - containerPort: 9000  # API
            - containerPort: 9001  # Console
          env:
            - name: MINIO_ROOT_USER
              value: minioadmin
            - name: MINIO_ROOT_PASSWORD
              value: minioadmin
          volumeMounts:
            - name: data
              mountPath: /data
      volumes:
        - name: data
          emptyDir: {}   # für dev reicht emptyDir, kein PVC nötig
---
apiVersion: v1
kind: Service
metadata:
  name: minio
  namespace: trip-manager-dev
spec:
  selector:
    app: minio
  ports:
    - name: api
      port: 9000
      targetPort: 9000
    - name: console
      port: 9001
      targetPort: 9001
EOF

kubectl apply -f /tmp/minio.yaml