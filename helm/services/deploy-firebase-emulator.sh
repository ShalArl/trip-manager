#!/bin/bash

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
            - containerPort: 9099   # Auth
            - containerPort: 8080   # Firestore
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