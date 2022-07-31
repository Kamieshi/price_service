package main

import (
	"net"

	rds "github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"priceService/internal/handlers"
	"priceService/internal/service"
	"priceService/protoc"
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
	rep := service.RedisListener{Client: client}

	grpcServer := grpc.NewServer()
	protoc.RegisterOwnPriceStreamServer(grpcServer, &handlers.PriceServerImplement{RedisListener: &rep})
	log.Info("gRPC server start")
	log.Info(grpcServer.Serve(listener))
	log.Info("gRPC server Stop")
}
