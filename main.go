package main

import (
	"context"
	"net"

	rds "github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/Kamieshi/price_service/internal/handlers"
	"github.com/Kamieshi/price_service/internal/service"
	"github.com/Kamieshi/price_service/protoc"
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
	ctx := context.Background()
	ch, err := rep.ListenChanel(ctx)
	if err != nil {
		log.WithError(err).Fatal()
	}
	grpcServer := grpc.NewServer()
	protoc.RegisterOwnPriceStreamServer(grpcServer, &handlers.PriceServerImplement{RedisListener: &rep})
	protoc.RegisterCommonPriceStreamServer(grpcServer, handlers.NewCommonPriceStreamServerImplement(ctx, ch))
	log.Info("gRPC server start")
	log.Info(grpcServer.Serve(listener))
	log.Info("gRPC server Stop")
}
