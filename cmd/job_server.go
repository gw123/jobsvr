package cmd

import (
	"net"
	"os"
	"os/signal"

	"github.com/gw123/jobsvr"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/gw123/glog"
	"github.com/gw123/jobsvr/apis"
	"github.com/gw123/jobsvr/internal/interceptors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// serverCmd represents the server command
var jobServerCmd = &cobra.Command{
	Use:   "job",
	Short: "job",
	Long:  `start job server`,
	Run: func(cmd *cobra.Command, args []string) {
		startJobServer()
	},
}

func init() {
	RootCmd.AddCommand(jobServerCmd)
}

func startJobServer() {
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

	jobService := apis.NewJobServer()
	jobsvr.RegisterJobManagerServer(server, jobService)

	go func() {
		if err := server.Serve(lis); err != nil {
			glog.DefaultLogger().Fatalf("grpc server err: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	server.Stop()
}
