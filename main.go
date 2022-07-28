package main

import (
	"context"
	"net"
	"runtime"

	"github.com/Kamieshi/price_service/internal/handler"
	"github.com/Kamieshi/price_service/internal/service"
	"github.com/Kamieshi/price_service/protoc"
	rds "github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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
	protoc.RegisterOwnPriceStreamServer(grpcServer, &handler.PriceServerImplement{RedisListener: &rep})
	protoc.RegisterCommonPriceStreamServer(grpcServer, handler.NewCommonPriceStreamServerImplement(ctx, ch))
	log.Info("gRPC server start")
	log.Info("Count core : ", runtime.NumCPU())
	log.Info(grpcServer.Serve(listener))
	log.Info("gRPC server Stop")
}
