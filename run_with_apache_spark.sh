#!/usr/bin/env bash

export BUCKET=BUCKER_NAME
export PROJECT=PROJECT_ID

go run main.go \
    --runner spark \
    --endpoint localhost:8099 \
    --file  gs://${BUCKET?}/amazon_reviews_sample_.csv \
    --output gs://${BUCKET?}/reporting_spark.txt \
    --project ${PROJECT?} \
    --temp_location gs://${BUCKET?}/tmp/ \
    --staging_location gs://${BUCKET?}/binaries/