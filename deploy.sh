#!/bin/bash

# Deploy script for Google Cloud Run
# Make sure you have gcloud CLI installed and configured

PROJECT_ID=${1:-"your-project-id"}
SERVICE_NAME="weather-api"
REGION="us-central1"

echo "Deploying to Google Cloud Run..."
echo "Project ID: $PROJECT_ID"
echo "Service Name: $SERVICE_NAME"
echo "Region: $REGION"

# Set the project
gcloud config set project $PROJECT_ID

# Enable required APIs
echo "Enabling required APIs..."
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com

# Deploy to Cloud Run
echo "Deploying to Cloud Run..."
gcloud run deploy $SERVICE_NAME \
  --source . \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --port 8080 \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 10

echo "Deployment complete!"
echo "Your service URL will be displayed above."
