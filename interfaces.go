package jobsvr

import (
	"context"

	"github.com/gw123/gworker"
)

type JobManager interface {
	Component
	GetJobManager() gworker.JobManager
}

type JobService interface {
	Component
	HandelJob(ctx context.Context, req *ListenQueueReq, server JobManager_ListenQueueServer) error
	PostJob(ctx context.Context, job *Job) error
}
