#!/bin/bash

set -ex

echo ${GoogleServiceAccountCredentials} >> key.json

gcloud auth activate-service-account concourse-deployment@selfhydro-197504.iam.gserviceaccount.com --key-file key.json --project selfhydro
gcloud functions deploy TransferStateToBigQuery --runtime go111 --env-vars-file .env.yaml --trigger-topic projects/selfhydro-197504/topics/telemetry-topic

rm key.json
