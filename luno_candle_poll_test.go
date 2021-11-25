package LunoCandlePoll

import (
	"context"
	"testing"
)


func TestLunoCandlePoll(t *testing.T) {
	tests := []struct { data string }{ {"Hello, World!\n"}}
	for _, test := range tests {
			m := PubSubMessage{
					Data: []byte(test.data),
			}
			err := LunoCandlePoll(context.Background(), m)

			if err != nil{
				t.Fatalf("LunoTickerPoll Failed: %v", err.Error())
				//t.Failed()
			}
	}
}