package service

import (
	"context"
	"encoding/json"

	"github.com/gw123/gworker/rabbiter"

	"github.com/gw123/jobsvr/constant"

	"github.com/gw123/jobsvr/di"

	"github.com/gw123/glog"
	"github.com/pkg/errors"

	"github.com/gw123/gworker"
	"github.com/gw123/jobsvr"
)

type JobService struct {
	jm gworker.JobManager
}

func NewJobService() *JobService {
	component, ok := di.GetComponent(constant.PkgJobManager)
	if !ok {
		glog.DefaultLogger().Fatal("component not register")
	}

	jm, ok := component.(jobsvr.JobManager)
	if !ok {
		glog.DefaultLogger().Fatal("component not implement *job_manager.JobManager")
	}

	return &JobService{
		jm: jm.GetJobManager(),
	}
}

func (j *JobService) Init() error {
	j = NewJobService()
	return nil
}

func (j JobService) Name() string {
	return constant.ServiceJob
}

func (j JobService) Close() error {
	return j.jm.Close()
}

func (j JobService) Reload() {
	j.Close()
	j.Init()
}

func (j JobService) GetWatchConfigKeys() jobsvr.WatchConfigKeys {
	return jobsvr.WatchConfigKeys{}
}

func (j JobService) OnConfigChange(key, val string) {
	return
}

func (j JobService) PostJob(ctx context.Context, job *jobsvr.Job) error {
	logger := glog.ExtractEntry(ctx).WithField("action", "PostJob")
	if err := j.jm.Dispatch(ctx, job); err != nil {
		logger.Errorf("dispatch err %v", err)
		return errors.Wrap(err, "Dispatch err")
	}
	return nil
}

func (j *JobService) HandelJob(ctx context.Context, req *jobsvr.ListenQueueReq, server jobsvr.JobManager_ListenQueueServer) error {
	queueName := req.GetName()
	glog.ExtractEntry(server.Context()).Infof("ListenQueue ...")
	var consumer rabbiter.Consumer
	var err error
	consumer, err = j.jm.Do(ctx, req.GetName(), func(ctx context.Context, jobber gworker.Jobber) error {
		var job jobsvr.Job
		if err := json.Unmarshal(jobber.Body(), &job); err != nil {
			return errors.Wrap(err, "json.Unmarshal Job")
		}

		jobStream := &jobsvr.JobStream{
			JobName: queueName,
			Length:  1,
			Job: []*jobsvr.Job{
				&job,
			},
		}

		if err := server.Send(jobStream); err != nil {
			glog.ExtractEntry(server.Context()).WithError(err).Errorf("send job to client stream err")
			consumer.Stop()
			return err
		}

		return nil
	})

	if err != nil {
		glog.ExtractEntry(ctx).WithError(err).Errorf("jm.Do")
	}
	consumer.Wait()
	return nil
}
