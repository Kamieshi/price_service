package handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"priceService/protoc"
)

func Loader(client protoc.CommonPriceStreamClient, countClients int, pr bool) {
	stream, err := client.GetPriceStream(context.Background(), &protoc.GetPriceStreamRequest{})
	if err != nil {
		log.WithError(err).Error()
	}
	count := 0
	sum := int64(0)
	for {
		data, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}

		tt, err := time.Parse("2006-01-02T15:04:05.000TZ-07:00", data.Time)

		sum += time.Since(tt).Nanoseconds()
		count += 1
		if count == 10 && pr {
			fmt.Println(countClients, "clients  AVG : ", time.Duration(sum/10))
			sum = 0
			count = 0
		}
	}
}

func TestHighLoadStreamPriceService(t *testing.T) {
	countClients := 1
	conn, err := grpc.Dial("localhost:5300", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client := protoc.NewCommonPriceStreamClient(conn)
	go Loader(client, countClients, true)
	for i := 0; i < countClients-1; i++ {
		go Loader(client, countClients, false)
	}
	time.Sleep(100 * time.Minute)
}
