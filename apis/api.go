package apis

import (
	"context"

	"github.com/gw123/jobsvr/constant"

	"github.com/gw123/jobsvr/di"

	"github.com/gw123/glog"
	"github.com/gw123/jobsvr"
	"github.com/pkg/errors"
)

type JobServer struct {
	jobSvr jobsvr.JobService
}

func NewJobServer() *JobServer {
	component, ok := di.GetComponent(constant.ServiceJob)
	if !ok {
		glog.DefaultLogger().Fatal("component not register")
	}

	jobSvr, ok := component.(jobsvr.JobService)
	if !ok {
		glog.DefaultLogger().Fatal("component not implement *job_manager.JobManager")
	}
	return &JobServer{
		jobSvr: jobSvr,
	}
}

func (j JobServer) PostJob(ctx context.Context, req *jobsvr.PostJobReq) (*jobsvr.PostJobResp, error) {

	if err := j.jobSvr.PostJob(ctx, req.GetJob()); err != nil {
		return nil, errors.Wrap(err, "jobSvr.PostJob")
	}

	resp := &jobsvr.PostJobResp{
		Uuid: req.GetJob().GetUuid(),
		Code: 200,
		Msg:  "success",
	}
	return resp, nil
}

func (j JobServer) PostJobs(ctx context.Context, req *jobsvr.PostJobsReq) (*jobsvr.PostJobsResp, error) {
	jobs := req.GetJobs()
	var errJobs []*jobsvr.Job
	for _, job := range jobs {
		if err := j.jobSvr.PostJob(ctx, job); err != nil {
			glog.ExtractEntry(ctx).WithError(err).Error("jobSvr.PostJob")
			errJobs = append(errJobs, job)
		}
	}

	resp := &jobsvr.PostJobsResp{
		Uuid:      req.GetUuid(),
		Code:      200,
		Msg:       "success",
		ErrorJobs: errJobs,
	}
	return resp, nil
}

func (j JobServer) ListenQueue(req *jobsvr.ListenQueueReq, server jobsvr.JobManager_ListenQueueServer) error {
	if req.GetSize() <= 0 {
		return errors.New("argument size err")
	}
	//size := req.GetSize()
	//queueName := req.GetName()
	glog.ExtractEntry(server.Context()).Infof("ListenQueue ...")
	err := j.jobSvr.HandelJob(server.Context(), req, server)
	if err != nil {
		glog.DefaultLogger().WithError(err).Error("jobSvr.HandelJob")
	}

	return err
}
