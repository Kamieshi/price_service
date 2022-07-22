package handlers

import (
	"context"
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGetPriceStream(t *testing.T) {
	conn, err := grpc.Dial("localhost:5300", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client := NewPriceClient(conn)
	stream, err := client.GetPriceStream(context.Background(), &GetPriceStreamRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for {
		data, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Println(data)
	}
}
