--DROP TABLE cryptobot-327815.luno.candles
CREATE TABLE cryptobot-327815.luno.candles (
    pair string
    ,open NUMERIC
    ,close NUMERIC
    ,high NUMERIC
    ,low NUMERIC
    ,volume NUMERIC
    ,timestamp TIMESTAMP
)
PARTITION BY DATE_TRUNC(timestamp, DAY)