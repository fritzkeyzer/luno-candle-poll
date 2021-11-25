#go mod init luno_candle_poll
#go mod tidy
#go mod vendor
#go build
#go test


#gcloud init
#gcloud functions deploy luno_candle_poll --entry-point LunoCandlePoll --trigger-topic LunoCandlePoll --region europe-west1 --runtime go116 --env-vars-file ./.env.yaml
#gcloud functions deploy luno_candles_poll --env-vars-file ./.env.yaml
#gcloud functions delete luno_candles_poll


#gcloud pubsub topics create LunoCandlesPoll
#gcloud pubsub topics publish LunoCandlesPoll --message LunoCandlesPoll
#gcloud pubsub topics delete LunoCandlesPoll


# TODO create cloud scheduler jobs
