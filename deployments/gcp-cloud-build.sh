#!/bin/bash

docker build -f ../.docker/Dockerfile --tag private-notes ../

docker image tag private-notes:latest private-notes:$1

# echo "Create service account"
# gcloud iam service-accounts create private-notes-sa --display-name="Private-Notes-SA" --project $1
# saEmail=$(gcloud iam service-accounts list --project $1 --format="value(email)" --filter="displayName=Private-Notes-SA")
# echo "Create the bucket"
# gcloud alpha storage buckets create gs://private-notes --project $1
# echo "Add $saEmail to gs storage"
# gcloud projects add-iam-policy-binding $1 \
#     --member=serviceAccount:$saEmail \
#     --role=roles/storage.objectAdmin \
#     --condition=title="private-notes-sa-binding",expression="resource.type == \"storage.googleapis.com/Bucket\" && resource.name == \"private-notes\""
# echo "Create service account key"
# if [ ! -f "key.json" ]; then
#    gcloud iam service-accounts keys create key.json --iam-account=$saEmail
# else
#     echo "Key already exists"
# fi

# cd ../

# gcloud functions deploy privateNotes --project $1 --runtime go116 --region=$2 --trigger-http --service-account $saEmail  --allow-unauthenticated --set-env-vars=GCP_PROJECT=$1,GCP_REGION=$2,GCP_CF_NAME=privateNotes
