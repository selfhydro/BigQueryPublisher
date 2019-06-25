# bigQueryPubliser

Cloud function that subscribes to GCP Pub/Sub topic and saves the received data into GCP Big Query. It subscribes to the telemetry-topic from [selfhydro](https://github.com/selfhydro/selfhydro)

### Deploy (first time)
`gcloud functions deploy TransferStateToBigQuery --runtime go111 --env-vars-file .env.yaml --trigger-topic telemetry-topic`
### Deploy
`gcloud functions deploy TransferStateToBigQuery --env-vars-file .env.yaml`
