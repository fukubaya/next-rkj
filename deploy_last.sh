#!/bin/sh

readonly REGION="asia-northeast1"
readonly TOPIC_LAST="daily-last"
readonly ENTORY_POINT_LAST="TweetLast"

gcloud functions deploy daily-last \
       --runtime go116 \
       --region "${REGION}" \
       --trigger-resource "${TOPIC_LAST}" \
       --trigger-event google.pubsub.topic.publish \
       --entry-point "${ENTORY_POINT_LAST}" \
       --env-vars-file .env.yaml
