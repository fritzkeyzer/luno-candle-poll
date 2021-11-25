#go mod init luno_candle_poll
#go mod tidy
#go mod vendor
#go build
#go test

# GOOGLE CLOUD REPO
#git add .
#git commit -a -m"Smiles"
#git push

# GOOGLE CLOUD FUNCTION (requires gcloud sdk to be installed)
#gcloud init
#gcloud functions deploy luno_candle_poll --entry-point LunoCandlePoll --trigger-topic LunoCandlePoll --region europe-west1 --runtime go116 --env-vars-file ./.env.yaml
#gcloud functions deploy luno_candles_poll --env-vars-file ./.env.yaml
#gcloud functions delete luno_candles_poll

# GOOGLE CLOUD PUB/SUB
#gcloud pubsub topics create LunoCandlesPoll
#gcloud pubsub topics publish LunoCandlesPoll --message LunoCandlesPoll
#gcloud pubsub topics delete LunoCandlesPoll


# TODO GOOGLE CLOUD SCHEDULER
