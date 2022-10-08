#!/bin/bash

cd ../

docker build --tag private-notes .

docker image tag private-notes:latest private-notes:$1

# # echo "Create service account"
gcloud iam service-accounts create private-notes-sa --display-name="Private-Notes-SA" --project $2
saEmail=$(gcloud iam service-accounts list --project $2 --format="value(email)" --filter="displayName=Private-Notes-SA")
# # echo "Create the bucket"
gcloud alpha storage buckets create gs://$3 --project $2 --location="$4"
# # echo "Add $saEmail to gs storage"
gcloud projects add-iam-policy-binding $2 \
    --member=serviceAccount:$saEmail \
    --role=roles/storage.objectAdmin \
    --condition=title="private-notes-sa-binding",expression="resource.type == \"storage.googleapis.com/Bucket\" && resource.name == \"$3\""
# # echo "Create service account key"
if [ ! -f "key.json" ]; then
   gcloud iam service-accounts keys create key.json --iam-account=$saEmail
else
    echo "Key already exists"
fi

gcloud artifacts repositories create private-notes --location=$4 --repository-format=docker --project $2 --description="Private Notes docker repository" 

gcloud builds submit --region=$4 --tag $4-docker.pkg.dev/$2/private-notes/private-notes:$1 .

gcloud run deploy private-notes --project $2 --service-account=$saEmail --region $4 --port 80 --image $4-docker.pkg.dev/$1/private-notes/private-notes:$1 --allow-unauthenticated --set-env-vars="ENVIRONMENT=cloudbuild" --set-env-vars="GCP_BUCKET_NAME=private-notes-zitec" --set-env-vars="PUBLIC_URL=private-notes-5wen2wh6lq-ey.a.run.app"
