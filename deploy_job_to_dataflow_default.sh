#!/usr/bin/env bash

export BUCKET=YOUR_BUCKET_NAME
export PROJECT=YOUR_RPOJECT_ID

go run main.go \
  --runner dataflow \
  --max_num_workers 10 \
  --job_name "demo-go" \
  --file gs://${BUCKET?}/amazon_reviews_sample*.csv \
  --output gs://${BUCKET?}/reporting.json \
  --project ${PROJECT?} \
  --temp_location gs://${BUCKET?}/tmp/ \
  --staging_location gs://${BUCKET?}/binaries/ \
  --worker_harness_container_image=YOUR_WORKER_HARNESS_CONTAINER_IMAGE
# image : gobeam-docker-techlead.bintray.io/beam/go:latest
