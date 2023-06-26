#!/bin/sh

readonly TOPIC="daily-topic"
readonly REGION="asia-northeast1"
readonly ENTORY_POINT="Tweet"

readonly TOPIC_SONG="daily-topic-song"
readonly ENTORY_POINT_SONG="TweetSong"

readonly TOPIC_YOUTUBE="daily-youtube"
readonly ENTORY_POINT_YOUTUBE="TweetYouTubeChannel"

readonly TOPIC_BOLT897="daily-bolt897"
readonly ENTORY_POINT_BOLT897="TweetBolt897"

gcloud functions deploy daily-tweet \
       --runtime go120 \
       --region "${REGION}" \
       --trigger-resource "${TOPIC}" \
       --trigger-event google.pubsub.topic.publish \
       --entry-point "${ENTORY_POINT}" \
       --env-vars-file .env.yaml

gcloud functions deploy daily-tweet-song \
       --runtime go120 \
       --region "${REGION}" \
       --trigger-resource "${TOPIC_SONG}" \
       --trigger-event google.pubsub.topic.publish \
       --entry-point "${ENTORY_POINT_SONG}" \
       --env-vars-file .env.yaml

gcloud functions deploy daily-youtube \
       --runtime go120 \
       --region "${REGION}" \
       --trigger-resource "${TOPIC_YOUTUBE}" \
       --trigger-event google.pubsub.topic.publish \
       --entry-point "${ENTORY_POINT_YOUTUBE}" \
       --env-vars-file .env.yaml

gcloud functions deploy daily-bolt897 \
       --runtime go120 \
       --region "${REGION}" \
       --trigger-resource "${TOPIC_BOLT897}" \
       --trigger-event google.pubsub.topic.publish \
       --entry-point "${ENTORY_POINT_BOLT897}" \
       --env-vars-file .env.yaml
