package main

import (
	"net"

	"github.com/gw123/jobsvr/apis"

	"github.com/gw123/jobsvr"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/gw123/jobsvr/internal/interceptors"
	"google.golang.org/grpc"

	"github.com/gw123/glog"
	"github.com/spf13/viper"
)

func main() {
	rpcPort := viper.GetString("rpc_port")
	if rpcPort == "" {
		rpcPort = "8082"
	}

	address := "0.0.0.0:" + rpcPort

	lis, err := net.Listen("tcp", address)
	if err != nil {
		glog.DefaultLogger().Fatalf("grpc listen err: %v", err)
	}

	glog.DefaultLogger().Infof("grpc listen success, address: %v", address)

	server := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			interceptors.Logger(glog.DefaultLogger())),
	)

	jobsvr.RegisterJobManagerServer(server, &apis.JobServer{})

	if err := server.Serve(lis); err != nil {
		glog.DefaultLogger().Fatalf("grpc server err: %v", err)
	}

	server.Stop()
}
