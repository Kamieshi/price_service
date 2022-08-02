package main

import (
	"fmt"
	"net"
	"runtime"

	"github.com/Kamieshi/price_service/internal/config"
	"github.com/Kamieshi/price_service/internal/handler"
	"github.com/Kamieshi/price_service/internal/service"
	"github.com/Kamieshi/price_service/protoc"
	rds "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Parse config from OS ENV")
		return
	}

	client := rds.NewClient(&rds.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	listener, err := net.Listen("tcp", ":"+conf.PortRPC)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	rep := service.RedisListener{Client: client}
	grpcServer := grpc.NewServer()
	protoc.RegisterOwnPriceStreamServer(grpcServer, &handler.PriceServerImplement{RedisListener: &rep})
	logrus.Info("gRPC server start")
	logrus.Info("Count core : ", runtime.NumCPU())
	logrus.Info("Start Port: ", conf.PortRPC)
	logrus.Info(grpcServer.Serve(listener))
	logrus.Info("gRPC server Stop")
}
