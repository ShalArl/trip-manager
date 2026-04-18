#!/bin/bash

source ../.env

echo "Begin project setup for $APP_NAME with project ID $MY_PROJECT_ID"


gcloud config set project "$MY_PROJECT_ID"