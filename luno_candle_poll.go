package LunoCandlePoll

import (
	"context"
	"errors"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	luno "github.com/luno/luno-go"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

var (
	PROJECT_ID = os.Getenv("PROJECT_ID")
	LUNO_API_ID = os.Getenv("LUNO_API_ID")
	LUNO_API_SECRET = os.Getenv("LUNO_API_SECRET")
	DEBUG_STRING = os.Getenv("DEBUG")
	DEBUG = false

	PAIR = os.Getenv("PAIR")
	CANDLE_DURATION_STRING = os.Getenv("CANDLE_DURATION")
	CANDLE_DURATION int64
)

// ENTRY POINT
func LunoCandlePoll(ctx context.Context, m PubSubMessage) error {
	CANDLE_DURATION, err := strconv.ParseInt(CANDLE_DURATION_STRING, 10, 64)
	if (err != nil){
		log.Println("ERROR: strconv.ParseInt: " + err.Error())
		return err
	}

	candle, err := getLatestCandle(PAIR, CANDLE_DURATION)
	if (err != nil){
		log.Println("ERROR: LunoGetCandles('"+PAIR+"', "+CANDLE_DURATION_STRING+"): " + err.Error())
		return err
	}
	log.Println(candle)

	err2 := bqInsertCandle(candle)
	if (err2 != nil){
		log.Println("ERROR: bqInsertCandle(): " + err.Error())
		return err2
	}

	return nil
}

func getLatestCandle(pair string, duration int64)(luno.Candle, error){

	now := time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
	var sinceWhen = now - 5*(duration*1000)

	//log.Print(sinceWhen)

	candles, err := lunoGetCandles(pair, sinceWhen)
	if err != nil {
		log.Println("ERROR: lunoGetCandles: " + err.Error())
		return luno.Candle{}, err
	}

	//spew.Dump(candles)
	if len(candles.Candles) == 0{
		var errMsg = "len(candles.Candles) == 0"
		return luno.Candle{}, errors.New(errMsg)
	}

	return candles.Candles[len(candles.Candles)-1], nil
}

func lunoGetCandles(pair string, since int64) (luno.GetCandlesResponse, error){

	CANDLE_DURATION, err := strconv.ParseInt(CANDLE_DURATION_STRING, 10, 64)
	if err != nil{
		log.Println("ERROR: strconv.ParseInt: " + err.Error())
		return luno.GetCandlesResponse{}, err
	}


	cl := luno.NewClient()
	//cl.SetDebug(*DEBUG)
	err2 := cl.SetAuth(LUNO_API_ID, LUNO_API_SECRET)
	if err2 != nil {
		log.Println("ERROR: cl.SetAuth: " + err2.Error())
		return luno.GetCandlesResponse{}, err2
	}

	ctx := context.Background()

	req := luno.GetCandlesRequest{Pair: "XBTZAR", Since: since, Duration: CANDLE_DURATION}		// if no currency is specified, luno returns all
	res, err := cl.GetCandles(ctx, &req)
	if err != nil {
		log.Println("ERROR: cl.GetCandles: " + err.Error())
		return luno.GetCandlesResponse{}, err
	}

	return *res, nil
}

func bqInsertCandle(candle luno.Candle)(error){

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, PROJECT_ID)
	if err != nil {
		log.Println("ERROR: bigquery.NewClient: " + err.Error())
		return err
	}
	defer func(client *bigquery.Client) {
		err := client.Close()
		if err != nil {
			log.Fatal("wtf happened here?")
			os.Exit(0)
		}
	}(client)

	if strings.Contains(PROJECT_ID, " "){
		return errors.New("ERROR: PROJECT_ID environment contains illegal chars")
	}

	bqTime, err := strconv.ParseInt(candle.Timestamp.QueryValue(), 10, 64)
	if err != nil {
		log.Println("ERROR: strconv.ParseInt: " + err.Error())
		return err
	}


	q := client.Query( `INSERT INTO cryptobot-327815.luno.candles (pair, open, close, high, low, volume, timestamp)
	VALUES (@pair, @open, @close, @high, @low, @volume, TIMESTAMP_MILLIS( @timestamp ) )`)
	q.Parameters = []bigquery.QueryParameter{
		{
			Name:  "pair",
			Value: PAIR,
		},
		{
			Name:  "open",
			Value: new(big.Rat).SetFloat64(candle.Open.Float64()),
		},
		{
			Name:  "close",
			Value: new(big.Rat).SetFloat64(candle.Close.Float64()),
		},
		{
			Name:  "high",
			Value: new(big.Rat).SetFloat64(candle.High.Float64()),
		},
		{
			Name:  "low",
			Value: new(big.Rat).SetFloat64(candle.Low.Float64()),
		},
		{
			Name:  "volume",
			Value: new(big.Rat).SetFloat64(candle.Volume.Float64()),
		},
		{
			Name: "timestamp",
			Value: bqTime,
		},
	}
	job, err := q.Run(ctx)
	if err != nil {
		log.Println("ERROR: q.Run: " + err.Error())
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		log.Println("ERROR: job.Wait: " + err.Error())
		return err
	}
	if err := status.Err(); err != nil {
		log.Println("ERROR: status: " + err.Error())
		return err
	}
        
	return nil
}