#!/bin/sh

readonly REGION="asia-northeast1"
readonly TOPIC_TOUR="daily-tour"
readonly ENTORY_POINT_TOUR="TweetFirstTour"

gcloud beta functions deploy daily-tour \
       --runtime go111 \
       --region "${REGION}" \
       --trigger-resource "${TOPIC_TOUR}" \
       --trigger-event google.pubsub.topic.publish \
       --entry-point "${ENTORY_POINT_TOUR}" \
       --env-vars-file .env.yaml
