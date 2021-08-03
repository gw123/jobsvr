package apis

import (
	"context"
	"errors"
	"time"

	"github.com/gw123/glog"

	"github.com/google/uuid"

	"github.com/gw123/jobsvr"
)

type JobServer struct {
}

func (j JobServer) ListenQueue(req *jobsvr.ListenQueueReq, server jobsvr.JobManager_ListenQueueServer) error {
	if req.GetSize() <= 0 {
		return errors.New("argument size err")
	}
	size := req.GetSize()
	queueName := req.GetName()
	glog.ExtractEntry(server.Context()).Infof("ListenQueue ...")

	for {
		jobStream := &jobsvr.JobStream{
			JobName: queueName,
			Length:  size,
			Job: []*jobsvr.Job{
				&jobsvr.Job{
					Uuid:      uuid.New().String(),
					Delay:     0,
					Name:      queueName,
					Payload:   nil,
					Timestamp: time.Now().Unix(),
				},
			},
		}

		if err := server.Send(jobStream); err != nil {
			glog.ExtractEntry(server.Context()).Errorf("jobStream")
			return err
		}
		time.Sleep(time.Second)
	}

	return nil
}

func (j JobServer) PostJob(ctx context.Context, req *jobsvr.PostJobReq) (*jobsvr.PostJobResp, error) {
	job := &jobsvr.PostJobResp{}
	job.Uuid = "xxxx"
	return job, nil
}

func (j JobServer) PostJobs(ctx context.Context, req *jobsvr.PostJobsReq) (*jobsvr.PostJobsResp, error) {
	job := &jobsvr.PostJobsResp{}
	job.Uuid = "xxxx2"
	return job, nil
}
