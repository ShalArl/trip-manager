#!/bin/bash

kubectl create configmap firebase-config \
  --namespace trip-manager-dev \
  --from-literal=firebase.json='{
    "emulators": {
      "auth": {
        "host": "0.0.0.0",
        "port": 9099
      },
      "firestore": {
        "host": "0.0.0.0",
        "port": 8080
      },
      "hub": {
        "host": "0.0.0.0",
        "port": 4400
      },
      "ui": {
        "host": "0.0.0.0",
        "port": 4000
      },
      "logging": {
        "host": "0.0.0.0",
        "port": 4500
      }
    }
  }'

cat > /tmp/firebase-emulator.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: firebase-emulator
  namespace: trip-manager-dev
spec:
  replicas: 1
  selector:
    matchLabels:
      app: firebase-emulator
  template:
    metadata:
      labels:
        app: firebase-emulator
    spec:
      containers:
        - name: firebase-emulator
          image: andreysenov/firebase-tools:latest
          command:
            - firebase
            - emulators:start
            - --only=auth,firestore
            - --project=trip-manager-local
          ports:
            - containerPort: 9099
            - containerPort: 8080
          volumeMounts:
            - name: firebase-config
              mountPath: /home/node/firebase.json
              subPath: firebase.json
      volumes:
        - name: firebase-config
          configMap:
            name: firebase-config
---
apiVersion: v1
kind: Service
metadata:
  name: firebase-emulator
  namespace: trip-manager-dev
spec:
  selector:
    app: firebase-emulator
  ports:
    - name: auth
      port: 9099
      targetPort: 9099
    - name: firestore
      port: 8080
      targetPort: 8080
EOF

kubectl apply -f /tmp/firebase-emulator.yaml