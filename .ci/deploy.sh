#!/bin/bash

set +e +x

cd ./BigQueryPublisher
echo ${GoogleServiceAccountCredentials} >> key.json
unset GoogleServiceAccountCredentials

set -ex


gcloud auth activate-service-account concourse-deployment@selfhydro-197504.iam.gserviceaccount.com --key-file key.json --project selfhydro
gcloud functions deploy TransferStateToBigQuery --runtime go111 --env-vars-file .env.yaml --trigger-topic telemetry-topic

rm key.json
