#!/bin/sh

readonly TOPIC="daily-topic"
readonly REGION="asia-northeast1"
readonly ENTORY_POINT="Tweet"

readonly TOPIC_SONG="daily-topic-song"
readonly ENTORY_POINT_SONG="TweetSong"

gcloud beta functions deploy daily-tweet \
       --runtime go111 \
       --region "${REGION}" \
       --trigger-resource "${TOPIC}" \
       --trigger-event google.pubsub.topic.publish \
       --entry-point "${ENTORY_POINT}" \
       --env-vars-file .env.yaml


gcloud beta functions deploy daily-tweet-song \
       --runtime go111 \
       --region "${REGION}" \
       --trigger-resource "${TOPIC_SONG}" \
       --trigger-event google.pubsub.topic.publish \
       --entry-point "${ENTORY_POINT_SONG}" \
       --env-vars-file .env.yaml
