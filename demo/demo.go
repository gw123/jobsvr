package main

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"

	"github.com/gw123/glog"
	"github.com/gw123/jobsvr"
	"google.golang.org/grpc"
)

const HeaderTraceID = "trace-id"

func main() {

	serverAddr := "127.0.0.1:8082"
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		return
	}
	client := jobsvr.NewJobManagerClient(conn)

	md := metadata.New(map[string]string{
		"uuid":     uuid.New().String(),
		"trace-id": uuid.New().String(), // 这里需要小写大写的字母grpc会自动转为小写
	})
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)
	resp, err := client.PostJob(ctx, &jobsvr.PostJobReq{Job: &jobsvr.Job{
		Uuid:      "xxxxx",
		Delay:     10,
		Name:      "abc",
		Payload:   nil,
		Timestamp: 0,
	}})
	if err != nil {
		glog.DefaultLogger().Error(err)
		return
	}

	stream, err := client.ListenQueue(ctx, &jobsvr.ListenQueueReq{
		Uuid:      uuid.New().String(),
		Name:      "xx",
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
		glog.Infof("recv: %+v", jobStreamData.String())
	}

	glog.DefaultLogger().Info(resp.String())
}
