package interceptors

import (
	"context"
	"errors"

	"github.com/gw123/glog"

	"github.com/gw123/glog/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const HeaderTraceID = "trace-id"
const HeaderUUID = "uuid"

// 日志中间件
func Logger(logger common.Logger) grpc.UnaryServerInterceptor {
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if ctx == nil {
			return nil, errors.New("ctx is nil")
		}
		ctx = glog.ToContext(ctx, logger)
		glog.AddPathname(ctx, info.FullMethod)
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if tmp, ok := md[HeaderTraceID]; ok && len(tmp) > 0 {
				glog.AddTraceID(ctx, tmp[0])
			}
			glog.ExtractEntry(ctx).Infof("metadata: %+v", md)
			//if tmp, ok := md[HeaderUUID]; ok && len(tmp) > 0 {
			//	glog.AddUUID(ctx, tmp[0])
			//}
		}

		glog.ExtractEntry(ctx).Infof("logger middle request: %v", req)
		response, err := handler(ctx, req)
		if err == nil {
			glog.ExtractEntry(ctx).Infof("logger middle response: %+v", response)
		} else {
			glog.ExtractEntry(ctx).Infof("logger middle error: %+v", err)
		}
		return response, err
	}
	return interceptor
}
