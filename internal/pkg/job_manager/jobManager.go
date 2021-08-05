package job_manager

import (
	"github.com/gw123/glog"
	"github.com/gw123/gworker"
	"github.com/gw123/jobsvr"
	"github.com/gw123/jobsvr/constant"
	"github.com/spf13/viper"

	"github.com/gw123/gworker/job"
)

type JobManager struct {
	jm *job.JobManager
}

func NewJobManager() jobsvr.JobManager {
	username := viper.GetString("queue.username")
	password := viper.GetString("queue.password")
	host := viper.GetString("queue.host")
	port := viper.GetInt("queue.port")
	vhost := viper.GetString("queue.vhost")
	if port == 0 {
		port = 5672
	}
	if vhost == "" {
		vhost = "/"
	}

	opts := job.JobManagerOptions{username, password, host, port, vhost}
	glog.DefaultLogger().Infof("job options %+v", opts)
	return &JobManager{
		jm: job.NewJobManager(opts),
	}
}

func (j *JobManager) Init() error {
	username := viper.GetString("queue.username")
	password := viper.GetString("queue.password")
	host := viper.GetString("queue.host")
	port := viper.GetInt("queue.port")
	vhost := viper.GetString("queue.vhost")
	if port == 0 {
		port = 5672
	}
	if vhost == "" {
		vhost = "/"
	}

	opts := job.JobManagerOptions{username, password, host, port, vhost}
	j.jm = job.NewJobManager(opts)
	return nil
}

func (j *JobManager) Name() string {
	return constant.PkgJobManager
}

func (j *JobManager) Close() error {
	if j.jm != nil {
		j.jm.Close()
	}
	return nil
}

func (j JobManager) Reload() {
	j.Close()
	j.Init()
}

func (j JobManager) GetWatchConfigKeys() jobsvr.WatchConfigKeys {
	return jobsvr.WatchConfigKeys{
		NeedReloadKeys: []string{"queue.username", "queue.password", "queue.host", "queue.port", "queue.vhost"},
	}
}

func (j JobManager) OnConfigChange(key, val string) {
	j.Close()
	j.Init()
}

func (j JobManager) GetJobManager() gworker.JobManager {
	return j.jm
}
