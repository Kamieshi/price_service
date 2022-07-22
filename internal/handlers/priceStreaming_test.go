package handlers

import (
	"context"
	"fmt"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestHighLoadStreamPriceService(t *testing.T) {
	countClients := 5000

	conn, err := grpc.Dial("localhost:5300", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client := NewPriceStreamingClient(conn)
	runClient := func(numb int) {
		stream, err := client.HighLoadStream(context.Background(), &GetPriceStreamRequest{})
		if err != nil {
			t.Fatal(err)
		}
		count := 0
		sum := int64(0)
		for {
			data, err := stream.Recv()
			if err != nil {
				log.Fatal(err)
			}
			t_s := time.Now()
			tt, err := time.Parse("2006-01-02T15:04:05.000TZ-07:00", data.Time)
			sum += t_s.Sub(tt).Nanoseconds()
			count += 1
			if count == 100 && numb == 0 {
				fmt.Println(countClients, "clients  AVG : ", time.Duration(sum/100))
				sum = 0
				count = 0
			}
		}
	}
	for i := 0; i < countClients; i++ {
		go runClient(i)
	}
	time.Sleep(100 * time.Second)
}