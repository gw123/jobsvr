package cmd

import (
	"github.com/gw123/glog"
	"github.com/gw123/gutils/apollo"
	"github.com/gw123/jobsvr/constant"
	"github.com/gw123/jobsvr/di"
	"github.com/gw123/jobsvr/internal/pkg/job_manager"
	"github.com/gw123/jobsvr/internal/service"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "jobsvr",
	Short: "jobsvr for echo ,esay to develop go web application",
	Long:  `jobsvr for echo ,esay to develop go web application`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		glog.DefaultLogger().Fatalf("execute %+v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	var cfgFile string
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func initConfig() {
	apollo.InitApolloFromDefaultEnv(constant.APPNAME)
	// 按需注册 优先注册pkg 后注册service
	di.RegisterComponent(job_manager.NewJobManager())
	di.RegisterComponent(service.NewJobService())
}
