#!/bin/bash

# Fail on any error.
set -e

export HOST_PROJECT_ID="pm-singleproject-20"
export SERVICE_PROJECT_ID="pm-test-10-e90f"
export FOLDER_ID="816243559598"
export USER_EMAIL_ID="parasmamgain@google.com"

#Create a service account in the host project to which the permission will be assigned
gcloud iam service-accounts create iac-sa \
    --description="iac-sa" \
    --display-name="iac-sa" \
    --project=$HOST_PROJECT_ID

echo "============ Setting Up Permission for the service account in Host Project ============"

gcloud projects add-iam-policy-binding $HOST_PROJECT_ID \
    --member="serviceAccount:iac-sa@$HOST_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/compute.networkAdmin"

gcloud projects add-iam-policy-binding $HOST_PROJECT_ID \
    --member="serviceAccount:iac-sa@$HOST_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountAdmin"

gcloud projects add-iam-policy-binding $HOST_PROJECT_ID \
    --member="serviceAccount:iac-sa@$HOST_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/serviceusage.serviceUsageAdmin"

gcloud projects add-iam-policy-binding $HOST_PROJECT_ID \
    --member="serviceAccount:iac-sa@$HOST_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/resourcemanager.projectIamAdmin"

gcloud projects add-iam-policy-binding $HOST_PROJECT_ID \
    --member="serviceAccount:iac-sa@$HOST_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/compute.securityAdmin"

gcloud iam service-accounts add-iam-policy-binding \
    iac-sa@$HOST_PROJECT_ID.iam.gserviceaccount.com \
    --member="user:$USER_EMAIL_ID" \
    --role="roles/iam.serviceAccountUser" \
    --project=$HOST_PROJECT_ID

gcloud iam service-accounts add-iam-policy-binding \
    iac-sa@$HOST_PROJECT_ID.iam.gserviceaccount.com \
    --member="user:$USER_EMAIL_ID" \
    --role="roles/iam.serviceAccountTokenCreator" \
    --project=$HOST_PROJECT_ID

echo "=========================================================================================="

echo "============ Setting Up Permission for the service account in Service Project ============"

gcloud projects add-iam-policy-binding $SERVICE_PROJECT_ID \
    --member="serviceAccount:iac-sa@$HOST_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/cloudsql.admin" \
    --project=$SERVICE_PROJECT_ID

gcloud projects add-iam-policy-binding $SERVICE_PROJECT_ID \
    --member="serviceAccount:iac-sa@$HOST_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/compute.instanceAdmin" \
    --project=$SERVICE_PROJECT_ID

gcloud projects add-iam-policy-binding $SERVICE_PROJECT_ID \
    --member="serviceAccount:iac-sa@$HOST_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountAdmin" \
    --project=$SERVICE_PROJECT_ID

gcloud projects add-iam-policy-binding $SERVICE_PROJECT_ID \
    --member="serviceAccount:iac-sa@$HOST_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountUser" \
    --project=$SERVICE_PROJECT_ID

gcloud projects add-iam-policy-binding $SERVICE_PROJECT_ID \
    --member="serviceAccount:iac-sa@$HOST_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/serviceusage.serviceUsageAdmin" \
    --project=$SERVICE_PROJECT_ID

gcloud projects add-iam-policy-binding $SERVICE_PROJECT_ID \
    --member="serviceAccount:iac-sa@$HOST_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/resourcemanager.projectIamAdmin" \
    --project=$SERVICE_PROJECT_ID

echo "=========================================================================================="
