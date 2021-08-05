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
	md := metadata.New(map[string]string{
		"uuid":     uuid.New().String(),
		"trace-id": uuid.New().String(), // 这里需要小写大写的字母grpc会自动转为小写
	})
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)
	resp, err := client.PostJob(ctx, &jobsvr.PostJobReq{Job: &jobsvr.Job{
		Uuid:        uuid.New().String(),
		DelaySecond: 0,
		Name:        "my_job",
		Payload:     nil,
		Timestamp:   0,
	}})
	if err != nil {
		glog.DefaultLogger().Error(err)
		return
	}
	glog.DefaultLogger().Infof("resp : %+v", resp)
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
	serverAddr := "127.0.0.1:8082"
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
