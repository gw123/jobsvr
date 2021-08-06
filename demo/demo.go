package main

import (
	"context"
	"flag"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"

	"github.com/gw123/glog"
	"github.com/gw123/jobsvr"
	"google.golang.org/grpc"
)

const HeaderTraceID = "trace-id"

func actionSend(client jobsvr.JobManagerClient) {
	ctx := context.Background()

	md := metadata.New(map[string]string{
		"uuid":     uuid.New().String(),
		"trace-id": uuid.New().String(), // 这里需要小写大写的字母grpc会自动转为小写
	})
	ctx = metadata.NewOutgoingContext(ctx, md)
	var jobs []*jobsvr.Job
	for i := 0; i < 1000; i++ {

		jobs = append(jobs,
			&jobsvr.Job{
				Uuid:        uuid.New().String(),
				DelaySecond: 0,
				Name:        "my_job",
				Payload:     nil,
				Timestamp:   0,
			})

		if len(jobs) == 100 {
			resp, err := client.PostJobs(ctx, &jobsvr.PostJobsReq{Jobs: jobs})
			if err != nil {
				glog.DefaultLogger().Error(err)
				return
			}
			glog.DefaultLogger().Infof("resp : %+v", resp)
			jobs = jobs[:0]
		}
	}

}

func actionListen(client jobsvr.JobManagerClient) {
	md := metadata.New(map[string]string{
		"uuid":     uuid.New().String(),
		"trace-id": uuid.New().String(), // 这里需要小写大写的字母grpc会自动转为小写
	})
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.ListenQueue(ctx, &jobsvr.ListenQueueReq{
		Uuid:      uuid.New().String(),
		Name:      "my_job",
		Timestamp: 0,
		Size:      10,
	})
	if err != nil {
		glog.DefaultLogger().Error(err)
		return
	}

	for {
		jobStreamData, err := stream.Recv()
		if err != nil {
			glog.DefaultLogger().Error(err)
			return
		}
		glog.Infof("recv job: %+v", jobStreamData.String())
	}
}

func main() {
	//serverAddr := "127.0.0.1:8082"
	serverAddr := "sh2:8082"
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return
	}
	client := jobsvr.NewJobManagerClient(conn)
	action := flag.String("action", "send", "send | listen")
	flag.Parse()

	switch *action {
	case "send":
		actionSend(client)
	case "listen":
		actionListen(client)
	}
}
