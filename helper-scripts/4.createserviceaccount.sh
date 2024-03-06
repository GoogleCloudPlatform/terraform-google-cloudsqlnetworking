#!/bin/bash
# Copyright 2023 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# Fail on any error.
set -e

export CONSUMER_PROJECT_ID=""
#SA_PROJECT_ID is the GCP project where the service account will be created
export SA_PROJECT_ID=""
export SERVICE_PROJECT_ID=""
export USER_PROJECT_ID=""
export FOLDER_ID=""
#IMPERONSATING_MEMBER will be of format "user:xyz@...com","group:xyz@...com" or "serviceAccount:xyz@...com"
export IMPERONSATING_MEMBER=""
export SA_NAME=""

gcloud config set project $SA_PROJECT_ID
echo "============ Creating Service Account for the service account                 ============"

#Create a service account in the CONSUMER project to which the permission will be assigned
gcloud iam service-accounts create $SA_NAME \
    --description="Service Account to be Used for creating GCP resources" \
    --display-name=$SA_NAME \
    --project=$SA_PROJECT_ID

echo "=========================================================================================="

echo "============ Setting Up storage admin Permission for the service account      ============"

gcloud projects add-iam-policy-binding $SA_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/storage.objectAdmin" \
    --project=$SA_PROJECT_ID

echo "=========================================================================================="

echo "============ Setting Up Permission for the service account in CONSUMER Project ==============="

gcloud projects add-iam-policy-binding $CONSUMER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/compute.networkAdmin" \
    --project=$CONSUMER_PROJECT_ID

gcloud projects add-iam-policy-binding $CONSUMER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountAdmin" \
    --project=$CONSUMER_PROJECT_ID

gcloud projects add-iam-policy-binding $_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/serviceusage.serviceUsageAdmin" \
    --project=$CONSUMER_PROJECT_ID

gcloud projects add-iam-policy-binding $CONSUMER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/resourcemanager.projectIamAdmin" \
    --project=$CONSUMER_PROJECT_ID

gcloud projects add-iam-policy-binding $CONSUMER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/compute.securityAdmin" \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com"  \
    --project=$CONSUMER_PROJECT_ID
    
gcloud projects add-iam-policy-binding $CONSUMER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/compute.instanceAdmin" \
    --project=$CONSUMER_PROJECT_ID

# Following permission are assigned to the User who can then impersonate this service account

gcloud iam service-accounts add-iam-policy-binding \
    $SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com \
    --member="$IMPERONSATING_MEMBER" \
    --role="roles/iam.serviceAccountUser" \
    --project=$SA_PROJECT_ID

gcloud iam service-accounts add-iam-policy-binding \
    $SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com \
    --member="$IMPERONSATING_MEMBER" \
    --role="roles/iam.serviceAccountTokenCreator" \
    --project=$SA_PROJECT_ID

echo "=========================================================================================="

echo "============ Setting Up Permission for the service account in Producer Project ============"

gcloud projects add-iam-policy-binding $PRODUCER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/cloudsql.admin" \
    --project=$PRODUCER_PROJECT_ID

gcloud projects add-iam-policy-binding $PRODUCER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountAdmin" \
    --project=$PRODUCER_PROJECT_ID

gcloud projects add-iam-policy-binding $PRODUCER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountUser" \
    --project=$PRODUCER_PROJECT_ID

gcloud projects add-iam-policy-binding $SERVICE_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/serviceusage.serviceUsageAdmin" \
    --project=$SERVICE_PROJECT_ID

gcloud projects add-iam-policy-binding $PRODUCER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/resourcemanager.projectIamAdmin" \
    --project=$PRODUCER_PROJECT_ID

echo "=========================================================================================="

echo "============ Setting Up Permission for the service account in User Project ==============="

gcloud projects add-iam-policy-binding $USER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/compute.networkAdmin" \
    --project=$USER_PROJECT_ID

gcloud projects add-iam-policy-binding $USER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/compute.securityAdmin" \
    --project=$USER_PROJECT_ID

gcloud projects add-iam-policy-binding $USER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/compute.instanceAdmin" \
    --project=$USER_PROJECT_ID

gcloud projects add-iam-policy-binding $USER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountAdmin" \
    --project=$USER_PROJECT_ID

gcloud projects add-iam-policy-binding $USER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountUser" \
    --project=$USER_PROJECT_ID

gcloud projects add-iam-policy-binding $USER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/serviceusage.serviceUsageAdmin" \
    --project=$USER_PROJECT_ID

gcloud projects add-iam-policy-binding $USER_PROJECT_ID \
    --member="serviceAccount:$SA_NAME@$SA_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/resourcemanager.projectIamAdmin" \
    --project=$USER_PROJECT_ID

echo "=========================================================================================="