package main

import (
	"context"
	"net"

	rds "github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"priceService/internal/handlers"
	"priceService/internal/repository"
)

func main() {
	client := rds.NewClient(&rds.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	listener, err := net.Listen("tcp", ":5300")
	if err != nil {
		log.WithError(err).Fatal()
	}
	rep := repository.Redis{Client: client}
	dtChan, err := rep.ListenChanel(context.Background())
	if err != nil {
		log.WithError(err).Fatal()
	}
	grpcServer := grpc.NewServer()
	handlers.RegisterPriceServer(grpcServer, &handlers.PriceServerImplement{Rep: &rep})
	handlers.RegisterPriceStreamingServer(grpcServer, handlers.NewPriceStreamingServerImplement(dtChan))
	log.Info("gRPC server start")
	log.Info(grpcServer.Serve(listener))
	log.Info("gRPC server Stop")
}
