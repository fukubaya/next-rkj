#!/bin/sh

readonly TOPIC="daily-topic"
readonly REGION="asia-northeast1"
readonly ENTORY_POINT="Tweet"

gcloud beta functions deploy daily-tweet \
       --runtime go111 \
       --region "${REGION}" \
       --trigger-resource "${TOPIC}" \
       --trigger-event google.pubsub.topic.publish \
       --entry-point "${ENTORY_POINT}" \
       --env-vars-file .env.yaml
